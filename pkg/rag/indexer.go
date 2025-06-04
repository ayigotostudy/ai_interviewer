package rag

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/embedding/openai"

	ri "github.com/cloudwego/eino-ext/components/indexer/redis"
	"github.com/cloudwego/eino/schema"
	"github.com/redis/go-redis/v9"
)

var indexer *ri.Indexer

func initIndexer(ctx context.Context, emb *openai.Embedder) (err error) {

	client := redis.NewClient(
		&redis.Options{
			Addr:     address,
			Password: password,
		},
	)

	//createIndex(ctx, client)

	// create es indexer component
	indexer, err = ri.NewIndexer(ctx, &ri.IndexerConfig{
		Client:    client,
		KeyPrefix: keyPrefix,
		DocumentToHashes: func(ctx context.Context, doc *schema.Document) (*ri.Hashes, error) {
			f2v := map[string]ri.FieldValue{
				// write doc.Content to field "content"
				// write vector of doc.Content to field "vector"
				customContentFieldName: {
					Value:     doc.Content,
					EmbedKey:  customContentVectorFieldName,
					Stringify: nil,
				},
			}

			return &ri.Hashes{
				Key:         doc.ID + "_suffix",
				Field2Value: f2v,
			}, nil
		},
		BatchSize: 5,
		Embedding: emb,
	})
	return
}

func GetIndexer() *ri.Indexer {
	return indexer
}

func createIndex(ctx context.Context, client *redis.Client) {
	// below use FT.CREATE to create an index.
	// see: https://redis.io/docs/latest/commands/ft.create/

	// schemas should match DocumentToHashes configured in IndexerConfig.
	schemas := []*redis.FieldSchema{
		{
			FieldName: customContentFieldName,
			FieldType: redis.SearchFieldTypeText,
		},
		{
			FieldName: customContentVectorFieldName,
			FieldType: redis.SearchFieldTypeVector,
			VectorArgs: &redis.FTVectorArgs{
				// FLAT index: https://redis.io/docs/latest/develop/interact/search-and-query/advanced-concepts/vectors/#flat-index
				// Choose the FLAT index when you have small datasets (< 1M vectors) or when perfect search accuracy is more important than search latency.
				FlatOptions: &redis.FTFlatOptions{
					Type:           "FLOAT32", // BFLOAT16 / FLOAT16 / FLOAT32 / FLOAT64. BFLOAT16 and FLOAT16 require v2.10 or later.
					Dim:            dimension, // keeps same with dimensions of Embedding
					DistanceMetric: "COSINE",  // L2 / IP / COSINE
				},
				// HNSW index: https://redis.io/docs/latest/develop/interact/search-and-query/advanced-concepts/vectors/#hnsw-index
				// HNSW, or hierarchical navigable small world, is an approximate nearest neighbors algorithm that uses a multi-layered graph to make vector search more scalable.
				HNSWOptions: nil,
			},
		},
	}

	options := &redis.FTCreateOptions{
		OnHash: true,
		Prefix: []any{keyPrefix},
	}

	result, err := client.FTCreate(ctx, indexName, options, schemas...).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println(result) // OK
}
