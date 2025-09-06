package router

import (
	"ai_jianli_go/internal/middleware"

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
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"}
	corsConfig.ExposeHeaders = []string{"Content-Length", "Content-Type", "X-RateLimit-Limit", "X-RateLimit-Remaining", "X-RateLimit-Reset"}
	corsConfig.AllowCredentials = false
	r.Use(cors.New(corsConfig))

	// 初始化限流器（从配置文件读取配置）
	middleware.InitRateLimiters()

	// 启动限流器清理协程
	middleware.StartRateLimitCleanup()

	// API版本分组
	v1 := r.Group("/api/v1")

	// 为不同模块应用限流中间件
	resume(v1.Group("/resume", middleware.GeneralRateLimitMiddleware()))
	meeting(v1.Group("/meeting", middleware.GeneralRateLimitMiddleware()))
	user(v1.Group("/user", middleware.GeneralRateLimitMiddleware()))
	speech(v1.Group("/speech", middleware.SpeechRateLimitMiddleware()))
	wiki(v1.Group("/wiki", middleware.GeneralRateLimitMiddleware()))

	// 限流管理接口（仅管理员可访问）
	ratelimit(v1.Group("/ratelimit", middleware.GeneralRateLimitMiddleware()))

	return r
}
