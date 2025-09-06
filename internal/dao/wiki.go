package dao

import (
	"ai_jianli_go/types/model"
	"ai_jianli_go/types/req"

	"gorm.io/gorm"
)

type WikiDAO struct {
	db *gorm.DB
}

func (w *WikiDAO) Create(wiki *model.Wiki) error {
	return w.db.Create(wiki).Error
}

func (w *WikiDAO) GetWikiList(userId uint) ([]*model.Wiki, int64, error) {
	var wikis []*model.Wiki
	var total int64

	query := w.db.Model(&model.Wiki{})

	// 添加查询条件
	if userId > 0 {
		query = query.Where("user_id = ?", userId)
	}

	// 获取总数
	query.Count(&total)

	// 查询所有结果（不分页）
	err := query.Find(&wikis).Error

	return wikis, total, err
}

func (w *WikiDAO) GetListByParentId(userId uint, parentId uint) ([]*model.Wiki, error) {
	var wikis []*model.Wiki

	err := w.db.Where("parent_id = ? AND user_id = ?", parentId, userId).Find(&wikis).Error
	return wikis, err
}

func (w *WikiDAO) GetWiki(id uint, userId uint) (*model.Wiki, error) {
	var wiki model.Wiki

	err := w.db.Where("id = ? AND user_id = ?", id, userId).First(&wiki).Error
	return &wiki, err
}

func (w *WikiDAO) DeleteWiki(request *req.DeleteWikiRequest) int64 {
	result := w.db.Where("id = ? AND user_id = ?", request.ID, request.UserID).Delete(&model.Wiki{})
	if result.Error != nil {
		return 0
	}
	return result.RowsAffected
}

func (w *WikiDAO) UpdateWiki(request *req.UpdateWikiRequest) int64 {
	result := w.db.Model(&model.Wiki{}).Where("id = ?", request.ID).Updates(map[string]interface{}{
		"title":   request.Title,
		"content": request.Content,
	})

	if result.Error != nil {
		return 0
	}
	return result.RowsAffected
}

func NewWikiDAO(db *gorm.DB) *WikiDAO {
	return &WikiDAO{db: db}
}
