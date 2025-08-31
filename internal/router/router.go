package router

import (
	"github.com/gin-contrib/cors"
	_ "github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	r := gin.Default()

	// 使用请求日志中间件
	// r.Use(middleware.RequestLogger())
	// pprof.Register(r)

	// 配置CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"}
	config.ExposeHeaders = []string{"Content-Length", "Content-Type"}
	config.AllowCredentials = false
	r.Use(cors.New(config))

	// API版本分组
	v1 := r.Group("/api/v1")
	resume(v1.Group("/resume"))
	meeting(v1.Group("/meeting"))
	user(v1.Group("/user"))
	speech(v1.Group("/speech"))
	wiki(v1.Group("/wiki"))

	return r
}
