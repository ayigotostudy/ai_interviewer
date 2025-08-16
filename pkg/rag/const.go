package rag

const (
	keyPrefix                    = "eino_doc_customized:"  // keyPrefix should be the prefix of keys you write to redis and want to retrieve.
	indexName                    = "test_index_customized" // indexName should be used in redis retriever.
	customContentFieldName       = "content"
	customContentVectorFieldName = "vector"
	customExtraFieldName         = "extra_field_number"
)


