# Wiki 文档索引和检索系统

这是一个基于 eino 框架的文档索引和检索系统，使用 Redis 作为向量数据库，OpenAI 作为嵌入模型。

## 主要特性

- 🚀 基于 eino 框架的高性能文档索引
- 🔍 语义搜索和向量相似度检索
- 📚 支持标题和内容的双重向量化
- 🗄️ Redis 向量数据库支持
- 🤖 OpenAI 嵌入模型集成
- 🛡️ 完善的错误处理和参数验证
- 📄 **多格式文档支持** - PDF、Markdown、CSV、TXT、LOG等

## 支持的文档格式

| 格式 | 扩展名 | 描述 | 处理方式 |
|------|--------|------|----------|
| PDF | `.pdf` | PDF文档 | 使用eino PDF解析器，按页面分割 |
| Markdown | `.md`, `.markdown` | Markdown文档 | 移除标记，保留纯文本，按章节分割 |
| CSV | `.csv` | CSV表格 | 保留标题行，为数据行添加行号标识 |
| 文本 | `.txt` | 纯文本文档 | 直接处理，按句子边界分块 |
| 日志 | `.log` | 日志文件 | 保留日志行，支持结构化信息 |

## 系统架构

```
Wiki Model
├── 索引器 (Indexer)
│   ├── 文档向量化
│   ├── Redis 存储
│   └── 批量处理
├── 检索器 (Retriever)
│   ├── 语义搜索
│   ├── 向量相似度计算
│   └── 结果排序
├── 文档处理器
│   ├── PDF解析器 (eino)
│   ├── 文本处理器 (自定义)
│   ├── Markdown处理器
│   ├── CSV处理器
│   └── 日志处理器
└── 数据库
    ├── MySQL (文档元数据)
    └── Redis (向量索引)
```

## 核心组件

### 1. Wiki 模型

```go
type Wiki struct {
    gorm.Model
    Title    string // 文章标题
    Content  string // 文章内容
    Url      string // 文章存储路径
    Type     int    // 文章类型: 0-知识库, 1-文件夹, 2-文章
    ParentID uint   // 父文件夹ID
    WikiType int    // 文章类型: 1-文档, 2-md文字, 3-文章链接
    UserId   uint   // 用户ID
    RootId   uint   // 根文件夹ID
}
```

### 2. 索引管理

- **CreateIndex**: 创建 Redis 搜索索引
- **DeleteIndex**: 删除索引
- **GetIndexInfo**: 获取索引信息

### 3. 文档操作

- **Store**: 批量存储文档
- **StoreSingle**: 存储单个文档
- **Search**: 语义搜索
- **SearchWithFilter**: 带过滤条件的搜索

### 4. 多格式文档处理

- **PDF文档**: 使用eino PDF解析器，支持页面分割
- **Markdown**: 智能移除标记，保留结构化内容
- **CSV表格**: 保留标题，为数据行添加标识
- **文本文件**: 智能分块，保持语义完整性
- **日志文件**: 保留原始格式，支持结构化查询

## 使用方法

### 1. 初始化

```go
// 创建 Wiki 实例
wiki := &Wiki{
    UserId: 1,
    RootId: 1,
    Title: "示例文档",
    Content: "文档内容...",
}

// 创建索引
err := wiki.CreateIndex(ctx, redisClient)

// 初始化索引器和检索器
err = wiki.Init(ctx, redisClient, embedder)
```

### 2. 存储不同类型的文档

```go
// 存储PDF文档
docs, err := rag.GetLoader().Load(ctx, document.Source{URI: "document.pdf"})
err = wiki.Store(ctx, docs)

// 存储Markdown文档
docs, err := rag.ParseTextFile("readme.md")
err = wiki.Store(ctx, docs)

// 存储CSV文档
docs, err := rag.ParseTextFile("data.csv")
err = wiki.Store(ctx, docs)

// 存储文本文档
docs, err := rag.ParseTextFile("notes.txt")
err = wiki.Store(ctx, docs)
```

### 3. 搜索文档

```go
// 基本搜索
results, err := wiki.Search(ctx, "搜索关键词")

// 带过滤条件的搜索
filter := map[string]interface{}{
    "file_type": ".pdf",
    "category": "技术",
}
results, err := wiki.SearchWithFilter(ctx, "搜索关键词", filter)
```

## 配置说明

### Redis 配置

```go
redisClient := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "", // 如果有密码
    DB:       0,
})
```

### OpenAI 配置

```go
embedder, err := openai.NewEmbedder(ctx, &openai.EmbeddingConfig{
    APIKey: "your-api-key",
    Model:  "text-embedding-3-small",
})
```

### 索引配置

- **向量维度**: 2560 (text-embedding-3-small)
- **距离度量**: COSINE
- **索引类型**: FLAT (适合小数据集)
- **批处理大小**: 10

## 文档处理特性

### PDF处理
- 自动页面分割
- 保持原始格式
- 支持多语言
- 向量化存储

### Markdown处理
- 智能标记移除
- 保留标题结构
- 代码块处理
- 列表格式化

### CSV处理
- 标题行识别
- 数据行编号
- 表格结构保持
- 支持大文件

### 文本处理
- 智能分块
- 句子边界识别
- 编码自动检测
- 内存优化

## 性能优化

1. **批量处理**: 支持批量文档存储，减少网络开销
2. **向量缓存**: Redis 内存存储，快速向量检索
3. **索引优化**: 使用 FLAT 索引提供精确的向量搜索
4. **并发支持**: 支持并发文档存储和检索
5. **智能分块**: 根据文档类型自动选择最佳分块策略

## 错误处理

系统提供了完善的错误处理机制：

- 文件格式验证错误
- 解析器初始化错误
- 索引创建/删除错误
- 文档存储错误
- 搜索错误
- 网络连接错误

## 扩展功能

### 1. 自定义文档处理器

```go
// 实现自定义的文档处理器
func (w *Wiki) ProcessCustomDocument(filePath string) ([]*schema.Document, error) {
    // 自定义处理逻辑
    content := readCustomFile(filePath)
    chunks := customChunking(content)
    
    return createDocuments(chunks, filePath), nil
}
```

### 2. 文档格式检测

```go
// 自动检测文档格式
if rag.IsSupportedFormat(filename) {
    fileType := rag.GetFileTypeDescription(filename)
    fmt.Printf("检测到文件类型: %s\n", fileType)
}
```

### 3. 批量文档处理

```go
// 批量处理多种格式的文档
documents := []string{
    "report.pdf",
    "data.csv", 
    "readme.md",
    "notes.txt",
}

for _, doc := range documents {
    if rag.IsSupportedFormat(doc) {
        docs, err := processDocument(doc)
        if err == nil {
            wiki.Store(ctx, docs)
        }
    }
}
```

## 注意事项

1. **内存使用**: 向量数据会占用较多内存，注意 Redis 内存配置
2. **API 限制**: OpenAI API 有调用频率限制，注意控制并发
3. **索引重建**: 修改索引结构后需要重建索引
4. **数据一致性**: 确保数据库和向量索引的数据同步
5. **文件编码**: 文本文件建议使用UTF-8编码
6. **大文件处理**: 超大文件会自动分块处理

## 故障排除

### 常见问题

1. **索引创建失败**
   - 检查 Redis 连接
   - 确认用户权限
   - 查看错误日志

2. **搜索无结果**
   - 检查索引是否正确创建
   - 确认文档是否成功存储
   - 验证搜索关键词

3. **性能问题**
   - 调整批处理大小
   - 优化向量维度
   - 检查网络延迟

4. **文档解析失败**
   - 检查文件格式是否支持
   - 确认文件是否损坏
   - 查看文件编码

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目。

## 许可证

MIT License
