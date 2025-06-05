package component

import (
	"ai_jianli_go/config"
	"ai_jianli_go/types/model"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 注册mysql
func initMySQL() {
	conf := config.GetMySQLConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.User, conf.Pwd, conf.Host, conf.Port, conf.DbName)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}
	// 设置表的字符集为 utf8mb4
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci").AutoMigrate(&model.Meeting{})
	initModel()
}

func GetMySQLDB() *gorm.DB {
	return db
}

func initModel() {
	db.AutoMigrate(model.User{})
	db.AutoMigrate(model.Meeting{})
}
