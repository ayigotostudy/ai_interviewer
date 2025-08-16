package rag

const (
	keyPrefix                    = "eino_doc_customized:"  // keyPrefix should be the prefix of keys you write to redis and want to retrieve.
	indexName                    = "test_index_customized" // indexName should be used in redis retriever.
	customContentFieldName       = "content"
	customContentVectorFieldName = "vector"
	customExtraFieldName         = "extra_field_number"
)

var (
	// redis
	password = ""               // 本地Redis通常没有密码
	address  = "localhost:6379" // 使用本地Redis
)


const (
	
)