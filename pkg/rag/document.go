package rag

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/document/loader/file"
	"github.com/cloudwego/eino-ext/components/document/loader/url"
	"github.com/cloudwego/eino-ext/components/document/parser/pdf"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/schema"
)

// 支持的文档类型
var supportedExtensions = map[string]string{
	".pdf":      "PDF文档",
	".txt":      "文本文件",
	".md":       "Markdown文档",
	".markdown": "Markdown文档",
	".docx":     "Word文档",
	".html":     "HTML文档",
	".htm":      "HTML文档",
}

// DocumentLoader 文档加载器接口
type DocumentLoader interface {
	Load(ctx context.Context, source string) ([]*schema.Document, error)
	GetSupportedFormats() []string
	GetLoaderName() string
}

// MultiFormatLoader 多功能文档加载器
type MultiFormatLoader struct {
	loaders map[string]DocumentLoader
	// 缓存已创建的loader
	loaderCache map[string]DocumentLoader
	// URL加载器（特殊处理，不需要扩展名）
	urlLoader DocumentLoader
}

// 全局多格式加载器实例
var multiLoader *MultiFormatLoader

// 初始化多格式文档加载器
func initLoder(ctx context.Context) (err error) {
	multiLoader = &MultiFormatLoader{
		loaders:     make(map[string]DocumentLoader),
		loaderCache: make(map[string]DocumentLoader),
	}

	// 注册各种格式的加载器
	if err := multiLoader.registerLoaders(ctx); err != nil {
		return fmt.Errorf("failed to register loaders: %w", err)
	}

	return nil
}

// registerLoaders 注册各种格式的加载器
func (multiLoader *MultiFormatLoader) registerLoaders(ctx context.Context) error {
	// 1. 注册PDF加载器
	pdfLoader, err := newPDFLoader(ctx)
	if err != nil {
		return fmt.Errorf("failed to create PDF loader: %w", err)
	}
	multiLoader.registerLoader(".pdf", pdfLoader)

	// 2. 注册文本文件加载器
	textLoader := newTextFileLoader()
	multiLoader.registerLoader(".txt", textLoader)
	multiLoader.registerLoader(".md", textLoader)
	multiLoader.registerLoader(".markdown", textLoader)

	// 3. 注册Word文档加载器
	docxLoader := newDOCXLoader()
	multiLoader.registerLoader(".docx", docxLoader)

	// 4. 注册HTML文件加载器
	htmlLoader := newHTMLLoader()
	multiLoader.registerLoader(".html", htmlLoader)
	multiLoader.registerLoader(".htm", htmlLoader)

	// 5. 注册URL加载器
	urlLoader, err := newURLLoader(ctx)
	if err != nil {
		return fmt.Errorf("failed to create URL loader: %w", err)
	}
	// URL加载器不需要扩展名，直接通过isURL函数检测
	multiLoader.urlLoader = urlLoader

	return nil
}

// registerLoader 注册加载器
func (ml *MultiFormatLoader) registerLoader(extension string, loader DocumentLoader) {
	ml.loaders[extension] = loader
}

// LoadDocument 根据文件类型或URL自动选择加载器加载文档
func LoadDocument(ctx context.Context, source string) ([]*schema.Document, error) {
	if multiLoader == nil {
		return nil, fmt.Errorf("multi-format loader not initialized")
	}

	// 检查是否为URL
	if isURL(source) {
		return multiLoader.loadWithURLLoader(ctx, source)
	}

	// 检查是否为本地文件
	if isLocalFile(source) {
		return multiLoader.loadWithFileLoader(ctx, source)
	}

	return nil, fmt.Errorf("unsupported source: %s", source)
}

// loadWithFileLoader 使用文件加载器加载文档
func (ml *MultiFormatLoader) loadWithFileLoader(ctx context.Context, filePath string) ([]*schema.Document, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	// 查找对应的加载器
	loader, exists := ml.loaders[ext]
	if !exists {
		// 如果没有找到专门的加载器，使用通用文本加载器
		loader = ml.loaders[".txt"]
		if loader == nil {
			return nil, fmt.Errorf("no suitable loader found for extension: %s", ext)
		}
	}

	return loader.Load(ctx, filePath)
}

// loadWithURLLoader 使用URL加载器加载文档
func (ml *MultiFormatLoader) loadWithURLLoader(ctx context.Context, url string) ([]*schema.Document, error) {
	if ml.urlLoader == nil {
		return nil, fmt.Errorf("URL loader not available")
	}

	return ml.urlLoader.Load(ctx, url)
}

// isURL 检查是否为URL
func isURL(source string) bool {
	return strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://")
}

// isLocalFile 检查是否为本地文件
func isLocalFile(source string) bool {
	if _, err := os.Stat(source); err == nil {
		return true
	}
	return false
}

// ==================== 各种加载器实现 ====================

// PDF加载器
type pdfLoader struct {
	fileLoader *file.FileLoader
}

func newPDFLoader(ctx context.Context) (*pdfLoader, error) {
	pdfParser, err := pdf.NewPDFParser(ctx, &pdf.Config{
		ToPages: true,
	})
	if err != nil {
		return nil, err
	}

	fileLoader, err := file.NewFileLoader(ctx, &file.FileLoaderConfig{
		UseNameAsID: true,
		Parser:      pdfParser,
	})
	if err != nil {
		return nil, err
	}

	return &pdfLoader{fileLoader: fileLoader}, nil
}

func (pl *pdfLoader) Load(ctx context.Context, source string) ([]*schema.Document, error) {
	docs, err := pl.fileLoader.Load(ctx, document.Source{URI: source})
	if err != nil {
		return nil, fmt.Errorf("failed to load PDF file %s: %w", source, err)
	}

	// 为每个文档添加元数据
	for _, doc := range docs {
		if doc.MetaData == nil {
			doc.MetaData = make(map[string]any)
		}
		doc.MetaData["source"] = source
		doc.MetaData["source_type"] = "file"
		doc.MetaData["file_type"] = ".pdf"
		doc.MetaData["file_name"] = filepath.Base(source)
		doc.MetaData["loader"] = "PDF专用加载器"
	}

	return docs, nil
}

func (pl *pdfLoader) GetSupportedFormats() []string {
	return []string{".pdf"}
}

func (pl *pdfLoader) GetLoaderName() string {
	return "PDF专用加载器"
}

// 文本文件加载器
type textFileLoader struct{}

func newTextFileLoader() *textFileLoader {
	return &textFileLoader{}
}

func (tfl *textFileLoader) Load(ctx context.Context, source string) ([]*schema.Document, error) {
	content, err := readFileContent(source)
	if err != nil {
		return nil, fmt.Errorf("failed to read text file %s: %w", source, err)
	}

	ext := strings.ToLower(filepath.Ext(source))
	var processedContent string

	switch ext {
	case ".md", ".markdown":
		processedContent = processMarkdown(content)
	default:
		processedContent = content
	}

	chunks := chunkText(processedContent, 1000)
	documents := make([]*schema.Document, len(chunks))

	for i, chunk := range chunks {
		documents[i] = &schema.Document{
			ID:      fmt.Sprintf("text_%s_chunk_%d", filepath.Base(source), i),
			Content: chunk,
			MetaData: map[string]any{
				"source":       source,
				"source_type":  "file",
				"file_type":    ext,
				"file_name":    filepath.Base(source),
				"chunk_id":     i,
				"total_chunks": len(chunks),
				"loader":       "文本文件加载器",
			},
		}
	}

	return documents, nil
}

func (tfl *textFileLoader) GetSupportedFormats() []string {
	return []string{".txt", ".md", ".markdown"}
}

func (tfl *textFileLoader) GetLoaderName() string {
	return "文本文件加载器"
}

// Word文档加载器
type docxLoader struct{}

func newDOCXLoader() *docxLoader {
	return &docxLoader{}
}

func (dl *docxLoader) Load(ctx context.Context, source string) ([]*schema.Document, error) {
	// 这里应该使用实际的DOCX解析器
	// 为了示例，我们创建一些模拟的文档
	content, err := readFileContent(source)
	if err != nil {
		return nil, fmt.Errorf("failed to read DOCX file %s: %w", source, err)
	}

	// 简单的文本提取（实际应该使用专门的DOCX解析器）
	docs := []*schema.Document{
		{
			ID:      fmt.Sprintf("docx_%s", filepath.Base(source)),
			Content: content,
			MetaData: map[string]any{
				"source":      source,
				"source_type": "file",
				"file_type":   ".docx",
				"file_name":   filepath.Base(source),
				"loader":      "Word文档加载器",
			},
		},
	}

	return docs, nil
}

func (dl *docxLoader) GetSupportedFormats() []string {
	return []string{".docx"}
}

func (dl *docxLoader) GetLoaderName() string {
	return "Word文档加载器"
}

// HTML文件加载器
type htmlLoader struct{}

func newHTMLLoader() *htmlLoader {
	return &htmlLoader{}
}

func (hl *htmlLoader) Load(ctx context.Context, source string) ([]*schema.Document, error) {
	content, err := readFileContent(source)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTML file %s: %w", source, err)
	}

	processedContent := cleanHTML(content)
	chunks := chunkText(processedContent, 1000)

	documents := make([]*schema.Document, len(chunks))
	for i, chunk := range chunks {
		documents[i] = &schema.Document{
			ID:      fmt.Sprintf("html_%s_chunk_%d", filepath.Base(source), i),
			Content: chunk,
			MetaData: map[string]any{
				"source":       source,
				"source_type":  "file",
				"file_type":    filepath.Ext(source),
				"file_name":    filepath.Base(source),
				"chunk_id":     i,
				"total_chunks": len(chunks),
				"loader":       "HTML文件加载器",
			},
		}
	}

	return documents, nil
}

func (hl *htmlLoader) GetSupportedFormats() []string {
	return []string{".html", ".htm"}
}

func (hl *htmlLoader) GetLoaderName() string {
	return "HTML文件加载器"
}

// URL加载器
type urlLoader struct {
	loader *url.Loader
}

func newURLLoader(ctx context.Context) (*urlLoader, error) {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	loader, err := url.NewLoader(ctx, &url.LoaderConfig{
		Client: httpClient,
	})
	if err != nil {
		return nil, err
	}

	return &urlLoader{loader: loader}, nil
}

func (ul *urlLoader) Load(ctx context.Context, source string) ([]*schema.Document, error) {
	docs, err := ul.loader.Load(ctx, document.Source{URI: source})
	if err != nil {
		return nil, fmt.Errorf("failed to load URL %s: %w", source, err)
	}

	// 为每个文档添加元数据
	for _, doc := range docs {
		if doc.MetaData == nil {
			doc.MetaData = make(map[string]any)
		}
		doc.MetaData["source"] = source
		doc.MetaData["source_type"] = "url"
		doc.MetaData["loaded_at"] = time.Now().Format(time.RFC3339)
		doc.MetaData["loader"] = "URL加载器"
	}

	return docs, nil
}

func (ul *urlLoader) GetSupportedFormats() []string {
	return []string{"http://", "https://"}
}

func (ul *urlLoader) GetLoaderName() string {
	return "URL加载器"
}

// ==================== 辅助函数 ====================

// readFileContent 读取文件内容
func readFileContent(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// processMarkdown 处理Markdown内容
func processMarkdown(content string) string {
	lines := strings.Split(content, "\n")
	var processedLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 移除标题标记
		line = strings.TrimPrefix(line, "#")
		line = strings.TrimPrefix(line, "##")
		line = strings.TrimPrefix(line, "###")
		line = strings.TrimPrefix(line, "####")
		line = strings.TrimPrefix(line, "#####")
		line = strings.TrimPrefix(line, "######")

		// 移除列表标记
		line = strings.TrimPrefix(line, "-")
		line = strings.TrimPrefix(line, "*")
		line = strings.TrimPrefix(line, "+")

		// 移除代码块标记
		line = strings.TrimPrefix(line, "```")
		line = strings.TrimSuffix(line, "```")

		line = strings.TrimSpace(line)
		if line != "" {
			processedLines = append(processedLines, line)
		}
	}

	return strings.Join(processedLines, "\n")
}

// cleanHTML 清理HTML标签
func cleanHTML(content string) string {
	// 简单的HTML标签清理
	content = strings.ReplaceAll(content, "<script>", "")
	content = strings.ReplaceAll(content, "</script>", "")
	content = strings.ReplaceAll(content, "<style>", "")
	content = strings.ReplaceAll(content, "</style>", "")
	content = strings.ReplaceAll(content, "<nav>", "")
	content = strings.ReplaceAll(content, "</nav>", "")
	content = strings.ReplaceAll(content, "<header>", "")
	content = strings.ReplaceAll(content, "</header>", "")
	content = strings.ReplaceAll(content, "<footer>", "")
	content = strings.ReplaceAll(content, "</footer>", "")

	return content
}

// chunkText 将文本分块
func chunkText(text string, chunkSize int) []string {
	if len(text) <= chunkSize {
		return []string{text}
	}

	var chunks []string
	start := 0

	for start < len(text) {
		end := start + chunkSize
		if end > len(text) {
			end = len(text)
		}

		// 尝试在句子边界分割
		if end < len(text) {
			// 向后查找最近的句号、换行符或空格
			for end > start && end < len(text) {
				if text[end] == '.' || text[end] == '\n' || text[end] == ' ' {
					break
				}
				end--
			}
		}

		chunk := strings.TrimSpace(text[start:end])
		if chunk != "" {
			chunks = append(chunks, chunk)
		}

		start = end
	}

	return chunks
}

// ==================== 公共接口函数 ====================

// GetSupportedFormats 获取支持的文档格式
func GetSupportedFormats() map[string]string {
	return supportedExtensions
}

// IsSupportedFormat 检查文件格式是否支持
func IsSupportedFormat(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	_, supported := supportedExtensions[ext]
	return supported
}

// GetFileTypeDescription 获取文件类型描述
func GetFileTypeDescription(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	if desc, exists := supportedExtensions[ext]; exists {
		return desc
	}
	return "未知格式"
}

// IsURL 检查是否为URL
func IsURL(source string) bool {
	return isURL(source)
}

// GetMultiLoader 获取多格式加载器
func GetMultiLoader() *MultiFormatLoader {
	return multiLoader
}

// GetLoaderInfo 获取加载器信息
func GetLoaderInfo() map[string]string {
	if multiLoader == nil {
		return nil
	}

	info := make(map[string]string)
	// 添加文件格式加载器
	for ext, loader := range multiLoader.loaders {
		info[ext] = loader.GetLoaderName()
	}
	// 添加URL加载器
	if multiLoader.urlLoader != nil {
		info["URL"] = multiLoader.urlLoader.GetLoaderName()
	}
	return info
}

// GetLoader 获取特定格式的加载器（向后兼容）
func GetLoader() *file.FileLoader {
	if multiLoader != nil {
		if loader, exists := multiLoader.loaders[".pdf"]; exists {
			if pdfLoader, ok := loader.(*pdfLoader); ok {
				return pdfLoader.fileLoader
			}
		}
	}
	return nil
}
