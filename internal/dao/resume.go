package dao

import (
	"ai_jianli_go/types/model"

	"gorm.io/gorm"
)

// ResumeDAO 简历数据访问对象
type ResumeDAO struct {
	db *gorm.DB
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
func (dao *ResumeDAO) UpdateResume(resume *model.Resume) error {
	return dao.db.Save(resume).Error
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
