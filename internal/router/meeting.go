package router

import (
	"ai_jianli_go/internal/middleware"

	"github.com/gin-gonic/gin"
)

func meeting(r *gin.RouterGroup) {

	r.Use(middleware.Auth())
	r.POST("/upload-resume")
	r.POST("/text/chat")

}
