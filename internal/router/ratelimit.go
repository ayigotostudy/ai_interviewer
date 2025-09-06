package router

import (
	ratelimitController "ai_jianli_go/internal/controller/ratelimit"

	"github.com/gin-gonic/gin"
)

// ratelimit 注册限流管理相关路由
func ratelimit(r *gin.RouterGroup) {
	controller := ratelimitController.NewRateLimitController()

	// 限流统计和管理接口（仅管理员可访问）
	r.GET("/stats", controller.GetStats)
	r.GET("/top-limited", controller.GetTopLimited)
	r.GET("/key/:key", controller.GetKeyStats)
	r.GET("/health", controller.GetHealth)
}
