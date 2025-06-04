package router

import (
	"ai_jianli_go/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	v1.Use(middleware.Cors())

	resume(v1.Group("/resume"))
	meeting(v1.Group("/meeting"))
	user(v1.Group("/user"))

	return r
}
