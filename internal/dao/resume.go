package dao

import (
	"ai_jianli_go/types/model"
	"context"

	"gorm.io/gorm"
)

// ResumeDAO 简历数据访问对象
type ResumeDAO struct {
	db *gorm.DB
}

func (dao *ResumeDAO) GetResumeList(ctx context.Context, userID uint) ([]*model.Resume, error) {
	resumes := make([]*model.Resume, 0)
	err := dao.db.Where("user_id = ?", userID).Find(&resumes).Error
	return resumes, err
}

func (dao *ResumeDAO) GetResume(id uint) (*model.Resume, error) {
	resume := new(model.Resume)
	err := dao.db.Where("id = ?", id).First(resume).Error
	return resume, err
}

// NewResumeDAO 创建简历DAO实例
func NewResumeDAO(db *gorm.DB) *ResumeDAO {
	return &ResumeDAO{db: db}
}

// CreateResume 创建新简历
func (dao *ResumeDAO) CreateResume(resume *model.Resume) error {
	return dao.db.Create(resume).Error
}

// GetResumeByID 根据ID获取简历
func (dao *ResumeDAO) GetResumeByID(id uint) (*model.Resume, error) {
	var resume model.Resume
	err := dao.db.Where("id = ?", id).First(&resume).Error
	return &resume, err
}

// GetResumesByUserID 获取用户的所有简历
func (dao *ResumeDAO) GetResumesByUserID(userID uint) ([]*model.Resume, error) {
	var resumes []*model.Resume
	err := dao.db.Where("user_id = ?", userID).Find(&resumes).Error
	return resumes, err
}

// UpdateResume 更新简历信息
func (dao *ResumeDAO) UpdateResume(id int64, content string) error {
	return dao.db.Model(&model.Resume{}).Where("id = ?", id).Update("content", content).Error
}

// DeleteResume 删除简历
func (dao *ResumeDAO) DeleteResume(id uint) error {
	return dao.db.Delete(&model.Resume{}, id).Error
}

func (dao *ResumeDAO) GetTemplate(id uint) (*model.Template, error) {
	var template model.Template
	err := dao.db.Where("id = ?", id).First(&template).Error
	return &template, err
}

func (dao *ResumeDAO) GetResumeTemplateList() ([]*model.Template, error) {
	var templates []*model.Template
	err := dao.db.Find(&templates).Error
	return templates, err
}
