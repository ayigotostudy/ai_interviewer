package component

import "context"

func Init() {
	ctx := context.Background()
	var err error
	err = initEmbedding(ctx)
	if err != nil {
		panic(err)
	}

	err = initIndexer(ctx, GetEmbedding())
	if err != nil {
		panic(err)
	}
	err = initRetriever(ctx, GetEmbedding())
	if err != nil {
		panic(err)
	}
	err = initLoder(ctx)
	if err != nil {
		panic(err)
	}
}
