package wikiService

import (
	"ai_jianli_go/component"
	"ai_jianli_go/internal/dao"
	"ai_jianli_go/logs"
	"ai_jianli_go/pkg/rag"
	"ai_jianli_go/types/model"
	"ai_jianli_go/types/req"
	"ai_jianli_go/types/resp/common"
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

type WikiService struct {
	wikiDAO *dao.WikiDAO
}

func NewWikiService(wikiDAO *dao.WikiDAO) *WikiService {
	return &WikiService{wikiDAO: wikiDAO}
}

func (s *WikiService) CreateWiki(request *req.CreateWikiRequest) int64 {

	wiki := &model.Wiki{
		UserId:   request.UserID,
		Title:    request.Title,
		ParentID: request.ParentID,
		WikiType: request.WikiType,
		Type:     request.Type,
		RootId:   request.RootId,
		Url:      request.Url,
	}

	switch request.Type {
	case model.WikiTypeKnowledge:
		err := s.wikiDAO.Create(wiki)
		if err != nil {
			logs.SugarLogger.Errorf("创建知识库失败: %v", err)
			return common.CodeCreateWikiFailed
		}
		wiki.RootId = wiki.ID
		err = wiki.CreateIndex(context.Background(), component.GetRedisDB())
		if err != nil {
			logs.SugarLogger.Errorf("创建知识库索引失败: %v", err)
			return common.CodeCreateIndexFailed
		}

	case model.WikiTypeFolder:
		err := s.wikiDAO.Create(wiki)
		if err != nil {
			logs.SugarLogger.Errorf("创建知识库失败: %v", err)
			return common.CodeCreateWikiFailed
		}
	case model.WikiTypeArticle:
		err := s.wikiDAO.Create(wiki)
		if err != nil {
			logs.SugarLogger.Errorf("创建知识库失败: %v", err)
			return common.CodeCreateWikiFailed
		}
		err = wiki.Init(context.Background(), component.GetRedisDB(), rag.GetEmbedding())
		if err != nil {
			logs.SugarLogger.Errorf("初始化知识库失败: %v", err)
			return common.CodeCreateWikiFailed
		}

		// 根据文件类型选择不同的处理方式
		docs, err := s.loadDocuments(wiki.Url)
		if err != nil {
			logs.SugarLogger.Errorf("加载知识库失败: %v", err)
			return common.CodeCreateWikiFailed
		}

		err = wiki.Store(context.Background(), docs)
		if err != nil {
			logs.SugarLogger.Errorf("存储知识库失败: %v", err)
			return common.CodeCreateWikiFailed
		}
	}
	return common.CodeSuccess
}

// loadDocuments 根据文件类型加载文档
func (s *WikiService) loadDocuments(filePath string) ([]*schema.Document, error) {
	// 使用新的统一文档加载器
	docs, err := rag.LoadDocument(context.Background(), filePath)
	if err != nil {
		return nil, fmt.Errorf("加载文档失败: %w", err)
	}

	return docs, nil
}

func (s *WikiService) GetWikiList(request *req.GetWikiListRequest) ([]*model.Wiki, int64) {
	return s.wikiDAO.GetWikiList(request.UserID)
}

func (s *WikiService) GetWiki(request *req.GetWikiRequest) (*model.Wiki, int64) {
	return s.wikiDAO.GetWiki(request.ID, request.UserID)
}

func (s *WikiService) DeleteWiki(request *req.DeleteWikiRequest) int64 {
	return s.wikiDAO.DeleteWiki(request)
}


func (s *WikiService) Query(request *req.QueryWikiRequest) (string, int64) {
	wiki := &model.Wiki{
		UserId: request.UserID,
		RootId: request.RootId,
	}
	wiki.Init(context.Background(), component.GetRedisDB(), rag.GetEmbedding())

	docs, err := wiki.Search(context.Background(), request.Query)
	if err != nil {
		logs.SugarLogger.Errorf("查询知识库失败: %v", err)
		return "", common.CodeQueryWikiFailed
	}

	// 构建上下文
	contexts := ""
	if len(docs) > 0 {
		contextParts := make([]string, len(docs))
		for i, doc := range docs {
			contextParts[i] = fmt.Sprintf("文档片段[%d]:\n%s\n", i+1, doc.Content)
		}
		contexts = strings.Join(contextParts, "\n---\n")
	}

	model := component.GetAIComponent().GetChatModel("gpt-4o")
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage(
			"你是AI知识助手，需要根据用户问题以及知识库内容，回答用户问题。",
		),
		schema.UserMessage("用户问题：\n{input}"),
		schema.AssistantMessage(
			"知识库内容：\n{context}",
			[]schema.ToolCall{},
		),
	)
	prompt := map[string]any{
		"input":  request.Query,
		"context": contexts,
	}
	ctx := context.Background()
	messages, err := template.Format(ctx, prompt)
	if err != nil {
		logs.SugarLogger.Errorf("生成提示失败: %v", err)
		return "", common.CodeQueryWikiFailed
	}

	resp, err := model.Generate(ctx, messages)
	if err != nil {
		logs.SugarLogger.Errorf("生成知识库回答失败: %v", err)
		return "", common.CodeQueryWikiFailed
	}

	return resp.Content, common.CodeSuccess
}
