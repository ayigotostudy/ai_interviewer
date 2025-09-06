# AI简历助手系统

这是一个基于 Go 语言开发的智能简历管理系统，集成了文档索引检索、AI面试、语音识别、知识库管理等功能。系统采用微服务架构，支持多用户、多角色权限管理。

## 主要特性

- 🚀 **智能简历管理** - 支持多种简历模板，在线编辑和预览
- 🤖 **AI面试系统** - 基于大语言模型的智能面试对话
- 🎤 **语音识别** - 集成语音转文字功能，支持面试录音
- 📚 **知识库管理** - 多格式文档索引和语义搜索
- 🔍 **智能检索** - 基于向量的语义搜索和相似度匹配
- 👥 **多用户系统** - 支持用户注册、登录和权限管理
- 🛡️ **权限控制** - 基于Casbin的细粒度权限管理
- 📄 **多格式支持** - PDF、Markdown、CSV、TXT、DOCX等文档处理
- 🗄️ **数据存储** - MySQL + Redis 双数据库架构
- 🌐 **RESTful API** - 完整的REST API接口

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
AI简历助手系统
├── 前端层 (Web Frontend)
│   ├── Vue.js 用户界面
│   ├── 简历编辑器
│   ├── 面试界面
│   └── 管理后台
├── API网关层
│   ├── Gin Web框架
│   ├── JWT认证中间件
│   ├── CORS跨域处理
│   └── 请求日志记录
├── 业务服务层
│   ├── 用户管理服务
│   ├── 简历管理服务
│   ├── 面试管理服务
│   ├── 知识库服务
│   └── 语音识别服务
├── AI服务层
│   ├── OpenAI集成
│   ├── 文档向量化
│   ├── 语义搜索
│   └── 智能问答
├── 数据访问层
│   ├── GORM ORM
│   ├── MySQL数据库
│   ├── Redis缓存
│   └── 文件存储
└── 基础设施层
    ├── 权限管理 (Casbin)
    ├── 日志系统
    ├── 配置管理
    └── 错误处理
```

## 功能模块详解

### 1. 用户管理模块 (User Management)

**功能描述**: 提供用户注册、登录、认证和权限管理功能

**核心特性**:
- 用户注册和登录
- JWT Token认证
- 密码加密存储
- 用户信息管理
- 角色权限控制

**API接口**:
- `POST /api/v1/user/register` - 用户注册
- `POST /api/v1/user/login` - 用户登录

**技术实现**:
- 使用Gin框架处理HTTP请求
- JWT Token进行身份认证
- Bcrypt加密用户密码
- GORM进行数据库操作

### 2. 简历管理模块 (Resume Management)

**功能描述**: 提供简历的创建、编辑、查看、删除和模板管理功能

**核心特性**:
- 多种简历模板支持
- 在线简历编辑器
- 简历预览功能
- 简历模板库
- 简历数据导出

**API接口**:
- `POST /api/v1/resume` - 创建简历
- `GET /api/v1/resume/list` - 获取简历列表
- `GET /api/v1/resume` - 获取单个简历
- `PUT /api/v1/resume` - 更新简历
- `DELETE /api/v1/resume` - 删除简历
- `GET /api/v1/resume/template` - 获取简历模板

**技术实现**:
- 支持JSON格式的简历数据结构
- 模板引擎渲染简历
- 文件上传和存储
- 数据验证和格式化

### 3. 面试管理模块 (Interview Management)

**功能描述**: 提供AI面试、面试记录、评价和简历上传功能

**核心特性**:
- AI智能面试对话
- 面试过程录制
- 面试评价生成
- 简历自动分析
- 面试历史记录

**API接口**:
- `POST /api/v1/meeting` - 创建面试
- `GET /api/v1/meeting` - 获取面试详情
- `PUT /api/v1/meeting` - 更新面试信息
- `DELETE /api/v1/meeting` - 删除面试
- `GET /api/v1/meeting/list` - 获取面试列表
- `POST /api/v1/meeting/upload_resume` - 上传简历
- `POST /api/v1/meeting/ai_interview` - AI面试对话
- `GET /api/v1/meeting/remark` - 获取面试评价

**技术实现**:
- 集成OpenAI GPT模型
- 实时对话处理
- 语音转文字功能
- 智能评价算法

### 4. 语音识别模块 (Speech Recognition)

**功能描述**: 提供语音转文字功能，支持面试录音和实时转换

**核心特性**:
- 实时语音识别
- 多语言支持
- 音频文件处理
- 识别结果优化
- 错误纠正机制

**API接口**:
- `POST /api/v1/speech/recognize` - 语音识别

**技术实现**:
- 集成科大讯飞语音识别API
- 支持多种音频格式
- 实时流式处理
- 识别结果后处理

### 5. 知识库管理模块 (Knowledge Base)

**功能描述**: 提供文档索引、存储、搜索和管理功能，支持多种文档格式

**核心特性**:
- 多格式文档支持 (PDF, MD, CSV, TXT, DOCX)
- 智能文档解析
- 向量化索引
- 语义搜索
- 文档分类管理

**API接口**:
- `POST /api/v1/wiki` - 创建知识库条目
- `GET /api/v1/wiki/list` - 获取知识库列表
- `GET /api/v1/wiki` - 获取单个条目
- `DELETE /api/v1/wiki` - 删除条目
- `POST /api/v1/wiki/query` - 语义搜索
- `GET /api/v1/wiki/file` - 获取文件
- `GET /api/v1/wiki/list/parent` - 按父级获取列表

**技术实现**:
- 基于eino框架的文档处理
- OpenAI嵌入模型向量化
- Redis向量数据库存储
- 智能文档分块和索引

### 6. 权限管理模块 (Permission Management)

**功能描述**: 基于Casbin的细粒度权限控制系统

**核心特性**:
- 多角色权限管理
- 资源访问控制
- 动态权限配置
- 权限继承机制
- 审计日志记录

**角色定义**:
- **游客 (Guest)**: 只能访问公开接口
- **普通用户 (Common)**: 拥有所有基础功能权限
- **会员 (Member)**: 高级功能权限
- **超级会员 (SuperMember)**: 批量处理权限
- **超级管理员 (SuperAdmin)**: 系统管理权限

**API接口**:
- `GET /api/v1/role/roles` - 获取角色列表
- `GET /api/v1/role/roles/:id` - 获取角色详情
- `GET /api/v1/role/permissions/check` - 权限检查

**技术实现**:
- Casbin权限控制引擎
- Redis策略存储
- 中间件权限验证
- 角色继承关系

## 技术栈

### 后端技术
- **Go 1.21+** - 主要编程语言
- **Gin** - Web框架
- **GORM** - ORM数据库操作
- **MySQL** - 关系型数据库
- **Redis** - 缓存和向量数据库
- **JWT** - 身份认证
- **Casbin** - 权限控制

### AI技术
- **OpenAI API** - 大语言模型和嵌入
- **科大讯飞** - 语音识别
- **eino框架** - 文档处理和向量化

### 前端技术
- **Vue.js** - 前端框架
- **Element UI** - UI组件库
- **Axios** - HTTP客户端

### 开发工具
- **Git** - 版本控制
- **Docker** - 容器化部署
- **Zap** - 日志系统
- **Viper** - 配置管理

## 快速开始

### 1. 环境要求

- Go 1.21+
- MySQL 8.0+
- Redis 6.0+
- Node.js 16+ (前端开发)

### 2. 安装依赖

```bash
# 克隆项目
git clone <repository-url>
cd ai_jianli_go

# 安装Go依赖
go mod tidy

# 安装前端依赖
cd web/aijianli-pro
npm install
```

### 3. 配置设置

#### 数据库配置
```yaml
# config/config.yaml
mysql:
  host: "localhost"
  port: "3306"
  user: "root"
  pwd: "your_password"
  dbname: "ai_jianli"

redis:
  host: "localhost"
  port: "6379"
  pwd: "your_redis_password"
```

#### AI服务配置
```yaml
# config/ai_config.yaml
openai:
  api_key: "your_openai_api_key"
  model: "gpt-3.5-turbo"

speech:
  api_key: "your_xunfei_api_key"
  api_secret: "your_xunfei_secret"
  app_id: "your_xunfei_app_id"
```

### 4. 数据库初始化

```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE ai_jianli CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 运行数据库迁移
go run main.go migrate
```

### 5. 启动服务

```bash
# 启动后端服务
go run main.go

# 启动前端服务 (开发环境)
cd web/aijianli-pro
npm run dev
```

### 6. API使用示例

#### 用户注册和登录
```bash
# 用户注册
curl -X POST http://localhost:8080/api/v1/user/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"123456","email":"test@example.com"}'

# 用户登录
curl -X POST http://localhost:8080/api/v1/user/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"123456"}'
```

#### 创建简历
```bash
curl -X POST http://localhost:8080/api/v1/resume \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "软件工程师简历",
    "content": "个人技能：Go, Python, JavaScript...",
    "template_id": 1
  }'
```

#### AI面试对话
```bash
curl -X POST http://localhost:8080/api/v1/meeting/ai_interview \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "meeting_id": 1,
    "question": "请介绍一下你的项目经验",
    "user_answer": "我参与过多个Web项目开发..."
  }'
```

#### 语音识别
```bash
curl -X POST http://localhost:8080/api/v1/speech/recognize \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "audio=@interview.wav"
```

#### 知识库搜索
```bash
curl -X POST http://localhost:8080/api/v1/wiki/query \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "query": "Go语言最佳实践",
    "user_id": 1,
    "root_id": 1
  }'
```

## 配置说明

### 系统配置

#### 数据库配置
```yaml
# config/config.yaml
mysql:
  host: "localhost"        # MySQL主机地址
  port: "3306"            # MySQL端口
  user: "root"            # 数据库用户名
  pwd: "your_password"    # 数据库密码
  dbname: "ai_jianli"     # 数据库名称

redis:
  host: "localhost"       # Redis主机地址
  port: "6379"           # Redis端口
  pwd: "your_password"   # Redis密码
```

#### 文件存储配置
```yaml
localPath:
  path: "/local/"         # 本地文件存储路径
```

#### 语音识别配置
```yaml
speech:
  apiKey: "your_api_key"     # 科大讯飞API Key
  apiSecret: "your_secret"   # 科大讯飞API Secret
  appId: "your_app_id"       # 科大讯飞应用ID
```

### AI服务配置

#### OpenAI配置
```yaml
# config/ai_config.yaml
openai:
  api_key: "sk-your-openai-api-key"
  model: "gpt-3.5-turbo"
  embedding_model: "text-embedding-3-small"
  max_tokens: 2000
  temperature: 0.7
```

#### 向量化配置
- **向量维度**: 2560 (text-embedding-3-small)
- **距离度量**: COSINE
- **索引类型**: FLAT
- **批处理大小**: 10

### 权限配置

#### Casbin配置
```yaml
role:
  model: "component/auth/casbin/model.conf"   # 权限模型文件
  policy: "component/auth/casbin/policy.csv"  # 权限策略文件
```

#### 角色权限
- **游客**: 只能访问登录注册接口
- **普通用户**: 拥有所有基础功能权限
- **会员**: 高级功能权限
- **超级会员**: 批量处理权限
- **超级管理员**: 系统管理权限

## 核心功能特性

### 智能简历管理
- **多模板支持**: 提供多种简历模板，满足不同行业需求
- **在线编辑**: 所见即所得的简历编辑器
- **实时预览**: 支持实时预览简历效果
- **数据导出**: 支持PDF、Word等格式导出
- **版本管理**: 简历版本历史记录和对比

### AI面试系统
- **智能问答**: 基于大语言模型的智能面试对话
- **个性化问题**: 根据简历内容生成个性化面试问题
- **实时评价**: 实时分析回答质量并给出建议
- **面试记录**: 完整的面试过程记录和回放
- **评分系统**: 多维度评分和综合评价

### 语音识别功能
- **实时转换**: 支持实时语音转文字
- **多语言支持**: 支持中文、英文等多种语言
- **音频处理**: 支持多种音频格式
- **噪音过滤**: 智能噪音过滤和语音增强
- **识别优化**: 针对面试场景的识别优化

### 知识库管理
- **多格式支持**: PDF、Markdown、CSV、TXT、DOCX等
- **智能解析**: 自动识别文档类型并解析内容
- **向量化索引**: 基于OpenAI嵌入的向量化存储
- **语义搜索**: 支持自然语言查询和语义匹配
- **分类管理**: 支持文档分类和标签管理

### 权限管理系统
- **角色管理**: 5级角色权限体系
- **细粒度控制**: 精确到接口级别的权限控制
- **动态配置**: 支持运行时权限配置
- **审计日志**: 完整的权限操作审计记录
- **继承机制**: 角色权限继承和覆盖机制

## 性能优化

### 系统性能
1. **数据库优化**: 使用GORM连接池，优化数据库查询
2. **缓存策略**: Redis缓存热点数据，减少数据库压力
3. **并发处理**: 支持高并发请求处理
4. **内存管理**: 智能内存分配和垃圾回收优化
5. **连接复用**: HTTP连接池和数据库连接复用

### AI服务优化
1. **批量处理**: 支持批量文档向量化，减少API调用
2. **向量缓存**: Redis存储向量数据，快速检索
3. **异步处理**: 异步处理耗时的AI任务
4. **请求合并**: 合并相似请求，减少重复计算
5. **智能分块**: 根据文档类型选择最佳分块策略

### 前端优化
1. **组件懒加载**: 按需加载组件，减少初始加载时间
2. **资源压缩**: 代码压缩和资源优化
3. **CDN加速**: 静态资源CDN分发
4. **缓存策略**: 浏览器缓存和本地存储
5. **响应式设计**: 适配不同设备和屏幕尺寸

## 部署指南

### Docker部署

#### 1. 构建镜像
```bash
# 构建后端镜像
docker build -t ai-jianli-backend .

# 构建前端镜像
cd web/aijianli-pro
docker build -t ai-jianli-frontend .
```

#### 2. Docker Compose部署
```yaml
# docker-compose.yml
version: '3.8'
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: your_password
      MYSQL_DATABASE: ai_jianli
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql

  redis:
    image: redis:6.0
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  backend:
    build: .
    ports:
      - "8080:8080"
    environment:
      - MYSQL_HOST=mysql
      - REDIS_HOST=redis
    depends_on:
      - mysql
      - redis

  frontend:
    build: ./web/aijianli-pro
    ports:
      - "3000:80"
    depends_on:
      - backend

volumes:
  mysql_data:
  redis_data:
```

### 生产环境部署

#### 1. 服务器要求
- **CPU**: 4核心以上
- **内存**: 8GB以上
- **存储**: 100GB以上SSD
- **网络**: 100Mbps以上带宽

#### 2. 环境配置
```bash
# 安装Go环境
wget https://golang.org/dl/go1.21.0.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz

# 安装MySQL
sudo apt update
sudo apt install mysql-server

# 安装Redis
sudo apt install redis-server

# 安装Nginx
sudo apt install nginx
```

#### 3. 服务配置
```bash
# 创建系统服务
sudo tee /etc/systemd/system/ai-jianli.service > /dev/null <<EOF
[Unit]
Description=AI简历助手后端服务
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/ai-jianli
ExecStart=/opt/ai-jianli/ai-jianli
Restart=always

[Install]
WantedBy=multi-user.target
EOF

# 启动服务
sudo systemctl enable ai-jianli
sudo systemctl start ai-jianli
```

### 监控和日志

#### 1. 日志管理
```bash
# 查看应用日志
tail -f logs/laboratory.log

# 查看错误日志
tail -f logs/laboratory.err.log

# 日志轮转配置
sudo tee /etc/logrotate.d/ai-jianli > /dev/null <<EOF
/opt/ai-jianli/logs/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 www-data www-data
}
EOF
```

#### 2. 性能监控
- 使用Prometheus + Grafana监控系统性能
- 监控数据库连接池状态
- 监控Redis缓存命中率
- 监控API响应时间

## 错误处理

系统提供了完善的错误处理机制：

### 业务错误
- 用户认证失败
- 权限不足错误
- 数据验证错误
- 业务逻辑错误

### 系统错误
- 数据库连接错误
- Redis连接错误
- 文件系统错误
- 网络连接错误

### AI服务错误
- OpenAI API调用失败
- 语音识别服务错误
- 向量化处理错误
- 模型推理错误

### 错误恢复
- 自动重试机制
- 降级服务策略
- 错误日志记录
- 异常告警通知

## 扩展功能

### 1. 自定义AI模型
```go
// 集成其他AI模型
type CustomAIService struct {
    APIKey string
    Model  string
}

func (c *CustomAIService) GenerateResponse(prompt string) (string, error) {
    // 自定义AI模型调用逻辑
    return response, nil
}
```

### 2. 插件系统
```go
// 插件接口定义
type Plugin interface {
    Name() string
    Version() string
    Initialize() error
    Execute(data interface{}) (interface{}, error)
}

// 注册插件
func RegisterPlugin(plugin Plugin) error {
    // 插件注册逻辑
    return nil
}
```

### 3. 多语言支持
```go
// 国际化支持
type I18nService struct {
    Locale string
    Messages map[string]string
}

func (i *I18nService) GetMessage(key string) string {
    // 多语言消息获取
    return i.Messages[key]
}
```

### 4. 数据分析
```go
// 用户行为分析
type AnalyticsService struct {
    db *gorm.DB
}

func (a *AnalyticsService) TrackUserAction(userID uint, action string) {
    // 用户行为追踪
}
```

## 开发指南

### 1. 代码规范
- 遵循Go官方代码规范
- 使用gofmt格式化代码
- 编写完整的单元测试
- 添加详细的注释和文档

### 2. 提交规范
```bash
# 提交信息格式
<type>(<scope>): <subject>

# 示例
feat(user): add user registration API
fix(auth): resolve JWT token validation issue
docs(readme): update installation guide
```

### 3. 测试规范
```bash
# 运行单元测试
go test ./...

# 运行集成测试
go test -tags=integration ./...

# 生成测试覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 常见问题

### 1. 安装问题
**Q: 如何解决依赖安装失败？**
A: 确保Go版本>=1.21，使用`go mod tidy`重新安装依赖

**Q: 数据库连接失败怎么办？**
A: 检查MySQL服务是否启动，确认配置文件中的数据库连接信息

### 2. 运行问题
**Q: 服务启动失败？**
A: 检查端口是否被占用，查看日志文件获取详细错误信息

**Q: AI功能无法使用？**
A: 确认OpenAI API密钥配置正确，检查网络连接

### 3. 性能问题
**Q: 系统响应慢？**
A: 检查数据库查询性能，优化Redis缓存配置

**Q: 内存使用过高？**
A: 调整向量批处理大小，优化文档分块策略

## 更新日志

### v1.0.0 (2024-01-01)
- 初始版本发布
- 基础简历管理功能
- AI面试系统
- 语音识别功能
- 知识库管理
- 权限管理系统

### 后续版本计划
- 移动端支持
- 更多AI模型集成
- 高级数据分析
- 企业版功能

## 贡献指南

### 1. 如何贡献
1. Fork 项目仓库
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

### 2. 贡献类型
- 🐛 Bug修复
- ✨ 新功能开发
- 📚 文档改进
- 🎨 代码优化
- ⚡ 性能提升

### 3. 开发环境
```bash
# 克隆项目
git clone https://github.com/your-username/ai_jianli_go.git

# 进入项目目录
cd ai_jianli_go

# 安装依赖
go mod tidy

# 运行测试
go test ./...

# 启动开发服务器
go run main.go
```

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 联系我们

- 项目主页: https://github.com/your-username/ai_jianli_go
- 问题反馈: https://github.com/your-username/ai_jianli_go/issues
- 邮箱: your-email@example.com

## 致谢

感谢以下开源项目的支持：
- [Gin](https://github.com/gin-gonic/gin) - Web框架
- [GORM](https://gorm.io/) - ORM库
- [Casbin](https://casbin.org/) - 权限控制
- [OpenAI](https://openai.com/) - AI服务
- [Vue.js](https://vuejs.org/) - 前端框架
