# API 接口更新总结

根据 `ai_go.md` 文件中的 API 接口定义，已对以下文件进行了更新：

## 已删除的接口

### 1. `src/service/interview.ts` - 完全删除
- 原因：`ai_go.md` 中没有 interview 相关的接口，只有 meeting 接口
- 功能重复：interview.ts 中的接口与 meeting.ts 重复

### 2. `src/service/user.ts` - 删除 `getUserInfo` 接口
- 原因：`ai_go.md` 中没有用户信息获取接口
- 保留接口：登录、注册、获取验证码

### 3. `src/service/resume.ts` - 删除多个接口
- 删除：获取列表、详情、更新、删除、上传文件
- 保留：创建简历接口
- 原因：`ai_go.md` 中只有创建简历的接口

## 已更新的接口

### 1. `src/service/meeting.ts` - 完全重构
- 添加了完整的 TypeScript 接口定义
- 修正了接口路径，确保与 `ai_go.md` 一致
- 更新了函数参数类型

### 2. 组件更新
- `src/views/resume/Index.vue` - 移除列表功能，只保留创建按钮
- `src/views/resume/Create.vue` - 更新参数格式，符合 API 定义
- `src/views/interview/Index.vue` - 使用 meeting 接口
- `src/views/interview/Detail.vue` - 使用 meeting 接口
- `src/views/interview/Room.vue` - 使用 meeting 接口

## 当前支持的 API 接口

### 用户模块
- `POST /user/login` - 用户登录
- `POST /user/register` - 用户注册
- `GET /authcode` - 获取验证码

### 简历模块
- `POST /resume` - 创建简历

### 面试模块
- `POST /meeting` - 创建面试
- `GET /meeting/list` - 获取面试列表
- `GET /meeting` - 获取面试详情
- `PUT /meeting` - 更新面试
- `DELETE /meeting` - 删除面试
- `POST /meeting/upload_resume` - 上传简历
- `POST /meeting/ai_interview` - AI 面试对话

### 语音识别
- `POST /api/v1/speech/recognize` - 语音识别

## 注意事项

1. **用户ID**: 在多个地方硬编码了 `user_id: 1`，实际使用时需要从用户状态获取
2. **模板ID**: 简历创建时硬编码了 `template_id: 1`，实际使用时需要选择合适的模板
3. **错误处理**: 所有 API 调用都包含了基本的错误处理
4. **类型安全**: 添加了完整的 TypeScript 接口定义，提高代码质量

## 下一步建议

1. 实现用户状态管理，获取真实的用户ID
2. 添加简历模板选择功能
3. 完善错误处理和用户提示
4. 添加加载状态和进度指示器
