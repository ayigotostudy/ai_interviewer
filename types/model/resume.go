package model

import (
	"gorm.io/gorm"
)

// Resume 简历模型
type Resume struct {
	gorm.Model 
	UserID     uint   `json:"user_id" gorm:"index"` // 关联用户ID
	Name       string `json:"name"`                  // 简历名称
	Content    string `json:"content"`               // 简历内容
	TemplateID int    `json:"template_id"`           // 模板ID
	Status     int    `json:"status"`                // 状态: 0-未完成, 1-已完成
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
	Name       string `json:"name"`
	Content    string `json:"content"`
}

// NewTemplate 创建模板实例
func NewTemplate(name string, content string) *Template {
	return &Template{
		Name:    name,
		Content: content,
	}
}