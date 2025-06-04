package component

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
	rr "github.com/cloudwego/eino-ext/components/retriever/redis"
	"github.com/cloudwego/eino/schema"
	"github.com/redis/go-redis/v9"
)

var retriever *rr.Retriever

func initRetriever(ctx context.Context, emb *openai.Embedder) (err error) {

	client := redis.NewClient(
		&redis.Options{
			Addr:          address,
			Password:      password,
			Protocol:      2,
			UnstableResp3: true,
		},
	)

	retriever, err = rr.NewRetriever(ctx, &rr.RetrieverConfig{
		Client:      client,
		Index:       indexName,
		VectorField: customContentVectorFieldName,
		Dialect:     2,
		ReturnFields: []string{
			customContentFieldName,
			customContentVectorFieldName,
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
				case customContentFieldName:
					resp.Content = v
				default:
					return nil, fmt.Errorf("unexpected field=%s", k)
				}
			}
			return resp, nil
		},
		TopK:      5,
		Embedding: emb,
	})
	if err != nil {
		panic(err)
	}

	return
}

func GetRetriever() *rr.Retriever {
	return retriever
}
