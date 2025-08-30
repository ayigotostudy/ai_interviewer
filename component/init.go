package component

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	db      *gorm.DB
	rediscc *redis.Client
)

func Init() {
	initMySQL()
	initRedis()
	initAIComponent()
}
