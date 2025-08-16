package router

import (
	speechController "ai_jianli_go/internal/controller/speech"

	"github.com/gin-gonic/gin"
)

// speech 注册语音识别相关路由
func speech(r *gin.RouterGroup) {
	controller := speechController.NewSpeechController()
	r.POST("/recognize", controller.Recognize)
}
