package model

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
	ri "github.com/cloudwego/eino-ext/components/indexer/redis"
	rr "github.com/cloudwego/eino-ext/components/retriever/redis"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	WikiTypeKnowledge = 0
	WikiTypeFolder    = 1
	WikiTypeArticle   = 2
)

type Wiki struct {
	gorm.Model
	Title    string `gorm:"column:title" json:"title"`         // 文章标题
	Url      string `gorm:"column:url" json:"url"`             // 文章存储路径
	Type     int    `gorm:"column:type" json:"type"`           // 文章类型, 0-知识库, 1-文件夹， 2-文章
	ParentID uint   `gorm:"column:parent_id" json:"parent_id"` // 父文件夹ID
	WikiType int    `gorm:"column:wiki_type" json:"wiki_type"` // 文章类型, 1-文档， 2-md文字， 3-文章链接
	UserId   uint   `gorm:"column:user_id" json:"user_id"`     // 用户ID
	RootId   uint   `gorm:"column:root_id" json:"root_id"`     // 根文件夹ID

	indexer   *ri.Indexer   `gorm:"-"`
	retriever *rr.Retriever `gorm:"-"`
}

func (w *Wiki) TableName() string {
	return "wiki"
}

func NewWiki() *Wiki {
	return &Wiki{}
}

const (
	wikiKeyPrefix                = "wiki_%d_%d:"      // keyPrefix: wiki_userId_rootId:docId
	wikiIndexName                = "index_wiki_%d_%d" // indexName: index_wiki_userId_rootId
	customContentFieldName       = "content"
	customTitleFieldName         = "title"
	customContentVectorFieldName = "vector"
	customTitleVectorFieldName   = "title_vector"
	customExtraFieldName         = "extra_field_number"
	dimension                    = 2560
)

// CreateIndex 创建Redis搜索索引
func (w *Wiki) CreateIndex(ctx context.Context, client *redis.Client) error {
	if w.UserId == 0 {
		return fmt.Errorf("invalid user_id")
	}

	keyPrefix := fmt.Sprintf(wikiKeyPrefix, w.UserId, w.RootId)
	indexName := fmt.Sprintf(wikiIndexName, w.UserId, w.RootId)

	// 检查索引是否已存在
	exists, err := client.Exists(ctx, indexName).Result()
	if err != nil {
		return fmt.Errorf("check index existence failed: %w", err)
	}
	if exists == 1 {
		// 索引已存在，删除旧索引
		_, err = client.FTDropIndex(ctx, indexName).Result()
		if err != nil {
			return fmt.Errorf("drop existing index failed: %w", err)
		}
	}

	// 创建新的搜索索引
	schemas := []*redis.FieldSchema{
		{
			FieldName: customContentFieldName,
			FieldType: redis.SearchFieldTypeText,
		},
		{
			FieldName: customTitleFieldName,
			FieldType: redis.SearchFieldTypeText,
		},
		{
			FieldName: customContentVectorFieldName,
			FieldType: redis.SearchFieldTypeVector,
			VectorArgs: &redis.FTVectorArgs{
				FlatOptions: &redis.FTFlatOptions{
					Type:           "FLOAT32",
					Dim:            dimension,
					DistanceMetric: "COSINE",
				},
			},
		},
		{
			FieldName: customTitleVectorFieldName,
			FieldType: redis.SearchFieldTypeVector,
			VectorArgs: &redis.FTVectorArgs{
				FlatOptions: &redis.FTFlatOptions{
					Type:           "FLOAT32",
					Dim:            dimension,
					DistanceMetric: "COSINE",
				},
			},
		},
	}

	options := &redis.FTCreateOptions{
		OnHash: true,
		Prefix: []any{keyPrefix},
	}

	_, err = client.FTCreate(ctx, indexName, options, schemas...).Result()
	if err != nil {
		return fmt.Errorf("create index failed: %w", err)
	}

	return nil
}

// Init 初始化索引器和检索器
func (w *Wiki) Init(ctx context.Context, client *redis.Client, emb *openai.Embedder) error {
	if w.UserId == 0 || w.RootId == 0 {
		return fmt.Errorf("invalid user_id or root_id")
	}

	keyPrefix := fmt.Sprintf(wikiKeyPrefix, w.UserId, w.RootId)
	indexName := fmt.Sprintf(wikiIndexName, w.UserId, w.RootId)

	// 创建索引器
	var err error
	w.indexer, err = ri.NewIndexer(ctx, &ri.IndexerConfig{
		Client:    client,
		KeyPrefix: keyPrefix,
		DocumentToHashes: func(ctx context.Context, doc *schema.Document) (*ri.Hashes, error) {
			// 生成唯一的文档ID
			docID := fmt.Sprintf("doc_%d_%d", w.ID, time.Now().UnixNano())

			f2v := map[string]ri.FieldValue{
				customContentFieldName: {
					Value:     doc.Content,
					EmbedKey:  customContentVectorFieldName,
					Stringify: nil,
				},
				customTitleFieldName: {
					Value:     doc.MetaData["title"].(string),
					EmbedKey:  customTitleVectorFieldName,
					Stringify: nil,
				},
			}

			return &ri.Hashes{
				Key:         docID,
				Field2Value: f2v,
			}, nil
		},
		BatchSize: 10, // 增加批处理大小
		Embedding: emb,
	})
	if err != nil {
		return fmt.Errorf("create indexer failed: %w", err)
	}

	// 创建检索器
	w.retriever, err = rr.NewRetriever(ctx, &rr.RetrieverConfig{
		Client:      client,
		Index:       indexName,
		VectorField: customContentVectorFieldName, // 主要基于内容向量搜索
		Dialect:     2,
		ReturnFields: []string{
			customContentFieldName,
			customTitleFieldName,
			customContentVectorFieldName,
			customTitleVectorFieldName,
		},
		DocumentConverter: func(ctx context.Context, doc redis.Document) (*schema.Document, error) {
			resp := &schema.Document{
				ID:       strings.TrimPrefix(doc.ID, keyPrefix),
				MetaData: map[string]any{},
			}

			for k, v := range doc.Fields {
				switch k {
				case customContentVectorFieldName:
					resp.WithDenseVector(rr.Bytes2Vector([]byte(v)))
				case customTitleVectorFieldName:
					// 标题向量作为元数据
					resp.MetaData["title_vector"] = rr.Bytes2Vector([]byte(v))
				case customContentFieldName:
					resp.Content = v
				case customTitleFieldName:
					resp.MetaData["title"] = v
				default:
					// 记录未知字段但不返回错误
					resp.MetaData[k] = v
				}
			}
			return resp, nil
		},
		TopK:      10, // 增加返回结果数量
		Embedding: emb,
	})
	if err != nil {
		return fmt.Errorf("create retriever failed: %w", err)
	}

	return nil
}

// Store 存储文档到索引
func (w *Wiki) Store(ctx context.Context, docs []*schema.Document, opts ...indexer.Option) error {
	if w.indexer == nil {
		return fmt.Errorf("indexer not initialized, call Init() first")
	}

	if len(docs) == 0 {
		return fmt.Errorf("no documents to store")
	}

	// 为每个文档添加必要的元数据
	for _, doc := range docs {
		if doc.MetaData == nil {
			doc.MetaData = make(map[string]any)
		}
		// 如果没有标题，使用内容的前50个字符作为标题
		if _, exists := doc.MetaData["title"]; !exists {
			title := doc.Content
			if len(title) > 50 {
				title = title[:50] + "..."
			}
			doc.MetaData["title"] = title
		}
	}

	_, err := w.indexer.Store(ctx, docs, opts...)
	if err != nil {
		return fmt.Errorf("store documents failed: %w", err)
	}

	return nil
}

// StoreSingle 存储单个文档
func (w *Wiki) StoreSingle(ctx context.Context, content, title string) error {
	doc := &schema.Document{
		Content: content,
		MetaData: map[string]any{
			"title": title,
		},
	}
	return w.Store(ctx, []*schema.Document{doc})
}

// Search 搜索文档
func (w *Wiki) Search(ctx context.Context, query string, opts ...retriever.Option) ([]*schema.Document, error) {
	if w.retriever == nil {
		return nil, fmt.Errorf("retriever not initialized, call Init() first")
	}

	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	docs, err := w.retriever.Retrieve(ctx, query, opts...)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	return docs, nil
}

// SearchWithFilter 带过滤条件的搜索
func (w *Wiki) SearchWithFilter(ctx context.Context, query string, filter map[string]interface{}, opts ...retriever.Option) ([]*schema.Document, error) {
	// 注意：eino框架可能不支持WithFilter，这里暂时跳过过滤条件
	// 可以根据实际需要实现自定义过滤逻辑

	return w.Search(ctx, query, opts...)
}

// DeleteIndex 删除索引
func (w *Wiki) DeleteIndex(ctx context.Context, client *redis.Client) error {
	if w.UserId == 0 || w.RootId == 0 {
		return fmt.Errorf("invalid user_id or root_id")
	}

	indexName := fmt.Sprintf(wikiIndexName, w.UserId, w.RootId)
	_, err := client.FTDropIndex(ctx, indexName).Result()
	if err != nil {
		return fmt.Errorf("delete index failed: %w", err)
	}

	return nil
}

// GetIndexInfo 获取索引信息
func (w *Wiki) GetIndexInfo(ctx context.Context, client *redis.Client) (redis.FTInfoResult, error) {
	if w.UserId == 0 || w.RootId == 0 {
		return redis.FTInfoResult{}, fmt.Errorf("invalid user_id or root_id")
	}

	indexName := fmt.Sprintf(wikiIndexName, w.UserId, w.RootId)
	info, err := client.FTInfo(ctx, indexName).Result()
	if err != nil {
		return redis.FTInfoResult{}, fmt.Errorf("get index info failed: %w", err)
	}

	return info, nil
}
