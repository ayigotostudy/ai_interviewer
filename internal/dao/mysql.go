package dao

import (
	"ai_jianli_go/logs"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMysql() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",
		"123456",
		"localhost",
		3306,
		"ai_jianli",
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logs.SugarLogger.Errorf("数据库连接失败: %v", err)
		panic(err)
	}
	return db
}
