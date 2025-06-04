package dao

import (
	"ai_jianli_go/types/model"

	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db: db}
}

func (dao *UserDAO) CreateUser(user *model.User) error {
	return dao.db.Create(user).Error
}

func (dao *UserDAO) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := dao.db.Where("email = ?", email).First(&user).Error
	return &user, err
}
