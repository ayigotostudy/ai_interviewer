package model

import (
	"ai_jianli_go/logs"
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"gorm.io/gorm"
)

// Resume 简历模型
type Resume struct {
	gorm.Model
	UserID     uint   `json:"user_id" gorm:"index"` // 关联用户ID
	Name       string `json:"name"`                 // 简历名称
	Content    string `json:"content"`              // 简历内容
	TemplateID int    `json:"template_id"`          // 模板ID
	Status     int    `json:"status"`               // 状态: 0-未完成, 1-已完成
}

// NewResume 创建简历实例
func NewResume(userID uint, name string, templateID int) *Resume {
	return &Resume{
		UserID:     userID,
		Name:       name,
		TemplateID: templateID,
		Status:     0, // 默认未完成
	}
}

// SetContent 设置简历内容
func (r *Resume) SetContent(content string) {
	r.Content = content
}

// Complete 标记简历为已完成
func (r *Resume) Complete() {
	r.Status = 1
}

type Template struct {
	gorm.Model
	Name        string `json:"name"`
	Content     string `json:"content"`
	ShowContent string `json:"show_content"` // 前端展示内容
}

// NewTemplate 创建模板实例
func NewTemplate(name string, content string) *Template {
	return &Template{
		Name:    name,
		Content: content,
	}
}

// SetShowContent 设置前端展示内容
func (r *Template) SetShowContent(model *openai.ChatModel) {
	if r.Content == "" {
		return
	}
	template := prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage("你是一个{role}。你需要根据用户的模版内容，生成数据来填充jinja模版, 返回数据为填充后的内容"),

		// 插入需要的对话历史（新对话的话这里不填）
		schema.MessagesPlaceholder("chat_history", true),

		// 用户消息模板
		schema.UserMessage("用户模版: {template}"),
	)
	ctx := context.Background()
	// 使用模板生成消息
	messages, err := template.Format(ctx, map[string]any{
		"role":     "模版填充数据专家",
		"template": r.Content,
	})
	if err != nil {
		logs.SugarLogger.Error("生成简历展示消息失败", err)
		return
	}
	result, err := model.Generate(ctx, messages)
	if err != nil {
		logs.SugarLogger.Error("生成简历展示消息失败", err)
		return
	}
	r.ShowContent = result.Content
}
