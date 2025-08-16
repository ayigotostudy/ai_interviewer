package component

import (
	"ai_jianli_go/config"

	"github.com/redis/go-redis/v9"
)

func GetRedisDB() *redis.Client {
	return rediscc
}

// 注册redis
func initRedis() {
	conf := config.GetRedisConfig()
	addr := conf.Host + ":" + conf.Port
	rediscc = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: conf.Pwd,
		DB:       0, // use default DB
		Protocol: 2,
	})
}
