package component

import (
	"ai_jianli_go/config"
	"ai_jianli_go/logs"
	"ai_jianli_go/types/model"
	"context"
	"fmt"
	"io"

	"os"

	"github.com/cloudwego/eino-ext/components/model/openai"
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
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci").AutoMigrate(&model.Meeting{}, &model.Resume{}, &model.Template{})
	initModel()
}

func GetMySQLDB() *gorm.DB {
	return db
}

func initModel() {
	db.AutoMigrate(model.User{})
	db.AutoMigrate(model.Meeting{})
	db.AutoMigrate(model.Resume{})
	db.AutoMigrate(model.Template{})
	// 初始化模板
	initTemplate()
}

func initTemplate() {
	db.AutoMigrate(model.Template{})
	data, _ := os.Getwd()
	file, _ := os.Open(data + "/template.md")
	content, _ := io.ReadAll(file)
	defer file.Close()
	// 初始化模板
	template := model.NewTemplate("默认模板", string(content))
	chatModel, err := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{
		Model:   "gpt-4o",
		BaseURL: "https://api.vveai.com/v1",
		APIKey:  "sk-npfmWk7VxIyeWYt23c5dCc49E7C343E487913c3e71E30b81",
	})
	if err != nil {
		logs.SugarLogger.Error("初始化chatModel失败", err)
		return
	}
	template.SetShowContent(chatModel)
	db.Create(template)
}
