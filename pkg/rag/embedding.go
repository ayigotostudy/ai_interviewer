package rag

import (
	"context"
	"time"

	"github.com/cloudwego/eino-ext/components/embedding/openai"
)

var (
	embedder *openai.Embedder
	// format    = aclopenai.EmbeddingEncodingFormatFloat
	dimension = 2560     // 1536 is the default dimension for text-embedding-ada-002
	user      = "system" // system is the default user for embedding

	embeddingAPIKey = "6dcde77d-fbc4-4ae4-9842-7f11a4ce61a1"
	Model           = "doubao-embedding-text-240715"
	BaseURL         = "https://ark.cn-beijing.volces.com/api/v3"
)

func initEmbedding(ctx context.Context) (err error) {
	embedder, err = openai.NewEmbedder(ctx, &openai.EmbeddingConfig{
		// OpenAI API 配置
		APIKey:  embeddingAPIKey,
		Model:   Model,
		Timeout: 30 * time.Second,

		BaseURL: BaseURL,

		Dimensions: &dimension, // 向量维度
		User:       &user,      // 用户标识
	})
	return
}

func GetEmbedding() *openai.Embedder {
	return embedder
}
