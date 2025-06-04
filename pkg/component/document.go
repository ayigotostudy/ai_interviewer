package component

import (
	"context"

	"github.com/cloudwego/eino-ext/components/document/loader/file"
	"github.com/cloudwego/eino-ext/components/document/parser/pdf"
)

var loader *file.FileLoader

func initLoder(ctx context.Context) (err error) {
	parser, err := pdf.NewPDFParser(ctx, &pdf.Config{
		ToPages: true, // 是否按页面分割文档
	})
	if err != nil {
		return err
	}
	loader, err = file.NewFileLoader(ctx, &file.FileLoaderConfig{
		UseNameAsID: true,   // 是否使用文件名作为文档ID
		Parser:      parser, // 可选：指定自定义解析器
	})
	return
}

func GetLoader() *file.FileLoader {
	return loader
}
