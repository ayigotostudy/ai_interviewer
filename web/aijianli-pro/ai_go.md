---
title: ai_go
language_tabs:
  - shell: Shell
  - http: HTTP
  - javascript: JavaScript
  - ruby: Ruby
  - python: Python
  - php: PHP
  - java: Java
  - go: Go
toc_footers: []
includes: []
search: true
code_clipboard: true
highlight_theme: darkula
headingLevel: 2
generator: "@tarslib/widdershins v4.0.30"

---

# ai_go

Base URLs:

# Authentication

- HTTP Authentication, scheme: bearer

# 用户基础模块

## POST 用户登录

POST /user/login

> Body 请求参数

```json
{
  "email": "2745969694@qq.com",
  "password": "666666"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{
  "code": "string",
  "msg": "string",
  "data": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|[基本响应](#schema基本响应)|

## POST 用户注册

POST /user/register

> Body 请求参数

```json
{
  "email": "2745969694@qq.com",
  "password": "666666",
  "auth_code": ""
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{
  "code": "string",
  "msg": "string",
  "data": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|[基本响应](#schema基本响应)|

## GET 用户获取验证码

GET /authcode

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|email|query|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "code": "string",
  "msg": "string",
  "data": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|[基本响应](#schema基本响应)|

# 简历生成模块

## POST 创建简历

POST /resume

> Body 请求参数

```json
{
  "user_id": 1,
  "name": "张三",
  "basic_info": "姓名：张三 | 年龄：28 | 学历：本科 | 联系方式：13800138000 | 邮箱：zhangsan@example.com",
  "work_exp": "1. 公司A - 软件工程师（2020-2022）：负责后端开发，优化系统性能。\n2. 公司B - 高级工程师（2022至今）：主导微服务架构升级。",
  "project_exp": "1. 智能客服系统（2021）：技术栈：Java/Spring Cloud，实现自动化响应。\n2. 数据中台建设（2023）：担任架构师，提升数据处理效率40%。",
  "self_eval": "全栈开发专家，擅长高并发系统设计，具备跨团队协作能力。",
  "awards": "2022年公司技术创新一等奖 | 2021年优秀员工",
  "target_job": "高级后端开发工程师，期望薪资25K-30K",
  "template_id": 4
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "msg": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» msg|string|true|none||none|

## GET 获取简历详情

GET /resume

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|id|query|string| 否 |none|

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "msg": "string",
  "data": {
    "user_id": 0,
    "name": "string",
    "content": "string",
    "template_id": 0,
    "status": 0
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» msg|string|true|none||none|
|» data|object|true|none||none|
|»» user_id|integer|true|none||none|
|»» name|string|true|none||none|
|»» content|string|true|none||none|
|»» template_id|integer|true|none||none|
|»» status|integer|true|none||none|

## DELETE 删除简历

DELETE /resume

> Body 请求参数

```json
{
  "id": 3
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## GET 获取简历模版

GET /resume/template

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "msg": "string",
  "data": [
    {
      "ID": 0,
      "CreatedAt": "string",
      "UpdatedAt": "string",
      "DeletedAt": null,
      "name": "string",
      "content": "string"
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» msg|string|true|none||none|
|» data|[object]|true|none||none|
|»» ID|integer|true|none||none|
|»» CreatedAt|string|true|none||none|
|»» UpdatedAt|string|true|none||none|
|»» DeletedAt|null|true|none||none|
|»» name|string|true|none||none|
|»» content|string|true|none||none|

## GET 获取我的简历列表

GET /resume/list

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "msg": "string",
  "data": [
    {
      "user_id": 0,
      "name": "string",
      "content": "string",
      "template_id": 0,
      "status": 0
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» msg|string|true|none||none|
|» data|[object]|true|none||none|
|»» user_id|integer|true|none||none|
|»» name|string|true|none||none|
|»» content|string|true|none||none|
|»» template_id|integer|true|none||none|
|»» status|integer|true|none||none|

# 面试模块

## POST 面试上传简历

POST /meeting/upload_resume

> Body 请求参数

```json
{
  "meeting_id": 2,
  "resume": "# 张三\n男 ｜ 28.000000岁 ｜ 高级后端开发工程师 ｜ 本科 ｜ 13800138000 ｜ zhangsan@example.com\n## 自我评价\n全栈开发专家，擅长高并发系统设计，具备跨团队协作能力。\n\n## 工作经历\n\n::: start\n**公司A**\n\n:::\n**软件工程师**\n\n:::\n**2020-2022**\n\n::: end\n负责后端开发，优化系统性能。\n\n::: start\n**公司B**\n\n:::\n**高级工程师**\n\n:::\n**2022至今**\n\n::: end\n主导微服务架构升级。\n\n\n## 项目经历\n\n::: start\n**智能客服系统**\n\n:::\n**开发人员**\n\n:::\n**Java/Spring Cloud**\n\n:::\n**2021**\n\n::: end\n实现自动化响应。\n\n::: start\n**数据中台建设**\n\n:::\n**架构师**\n\n:::\n**数据处理与架构设计**\n\n:::\n**2023**\n\n::: end\n提升数据处理效率40%。\n\n\n## 专业技能\nJava, Spring Cloud, 微服务架构, 数据处理, 高并发系统设计\n\n## 教育背景与资质\n::: start\n**某大学** | 计算机科学与技术 | 本科\n\n:::\n**2016-2020**\n\n::: end\n2022年公司技术创新一等奖, 2021年优秀员工"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{
  "code": "string",
  "msg": "string",
  "data": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|[基本响应](#schema基本响应)|

## POST 面试对话(文本)

POST /meeting/ai_interview

> Body 请求参数

```json
{
  "meeting_id": 2,
  "answer": "我认为在团队协作中，沟通透明度和目标对齐是关键。例如，在上一份工作中，我们通过每日站会同步进度，确保所有成员理解优先级..."
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{
  "code": "string",
  "msg": "string",
  "data": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|[基本响应](#schema基本响应)|

## POST 创建面试

POST /meeting

> Body 请求参数

```json
{
  "candidate": "张三",
  "position": "高级软件工程师",
  "job_description": "负责后端系统开发，要求精通Go语言和分布式系统",
  "time": 1735660800,
  "status": "planned",
  "remark": "候选人需提前准备算法题"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## PUT 更新面试内容

PUT /meeting

> Body 请求参数

```json
{
  "id": 1001,
  "user_id": 2025,
  "candidate": "李四",
  "position": "高级后端工程师",
  "job_description": "负责分布式系统设计与开发，要求精通Go语言和微服务架构",
  "time": 1734364800,
  "status": "in_progress",
  "remark": "候选人需准备过往项目案例",
  "interview_record": "Q: 请介绍你设计的分布式系统...\nA: 我采用分片策略解决数据一致性问题...",
  "interview_summary": "候选人技术扎实，但系统容错设计经验待考察"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## GET 获取面试详情

GET /meeting

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|id|query|string| 否 |none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

## DELETE 删除面试

DELETE /meeting

> Body 请求参数

```json
{
  "id": 1
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "msg": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» msg|string|true|none||none|

## GET 获取面试列表

GET /meeting/list

> 返回示例

> 200 Response

```json
{
  "code": 0,
  "msg": "string",
  "data": [
    {
      "id": 0,
      "created_at": "string",
      "updated_at": "string",
      "deleted_at": null,
      "user_id": 0,
      "candidate": "string",
      "position": "string",
      "job_description": "string",
      "time": 0,
      "status": "string",
      "remark": "string",
      "resume": "string",
      "interview_record": "string",
      "interview_summary": "string"
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» code|integer|true|none||none|
|» msg|string|true|none||none|
|» data|[object]|true|none||none|
|»» id|integer|false|none||none|
|»» created_at|string|false|none||none|
|»» updated_at|string|false|none||none|
|»» deleted_at|null|false|none||none|
|»» user_id|integer|false|none||none|
|»» candidate|string|false|none||none|
|»» position|string|false|none||none|
|»» job_description|string|false|none||none|
|»» time|integer|false|none||none|
|»» status|string|false|none||none|
|»» remark|string|false|none||none|
|»» resume|string|false|none||none|
|»» interview_record|string|false|none||none|
|»» interview_summary|string|false|none||none|

## GET 获取面试评价

GET /meeting/remark

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|id|query|string| 否 |none|

> 返回示例

> 200 Response

```json
{}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

# 公用模块

## POST 语音识别

POST /api/v1/speech/recognize

> Body 请求参数

```yaml
audio: ""

```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» audio|body|string(binary)| 否 |none|

> 返回示例

> 200 Response

```json
{
  "text": "string"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» text|string|true|none||none|

# 数据模型

<h2 id="tocS_基本响应">基本响应</h2>

<a id="schema基本响应"></a>
<a id="schema_基本响应"></a>
<a id="tocS基本响应"></a>
<a id="tocs基本响应"></a>

```json
{
  "code": "string",
  "msg": "string",
  "data": "string"
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|code|string|true|none||none|
|msg|string|true|none||none|
|data|any|true|none||none|

oneOf

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» *anonymous*|string|false|none||none|

xor

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» *anonymous*|integer|false|none||none|

xor

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» *anonymous*|boolean|false|none||none|

xor

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» *anonymous*|array|false|none||none|

xor

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» *anonymous*|object|false|none||none|

xor

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» *anonymous*|number|false|none||none|

