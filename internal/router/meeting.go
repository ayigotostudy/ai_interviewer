package router

import (
	"ai_jianli_go/component"
	meetingController "ai_jianli_go/internal/controller/meeting"
	"ai_jianli_go/internal/dao"
	"ai_jianli_go/internal/middleware"
	meetingService "ai_jianli_go/internal/service/meeting"

	"github.com/gin-gonic/gin"
)

func meeting(rg *gin.RouterGroup) {
	meetingDao := dao.NewMeetingDAO(component.GetMySQLDB())
	meetingSvc := meetingService.NewMeetingService(meetingDao)
	meetingCtrl := meetingController.NewMeetingController(meetingSvc)

	rg.Use(middleware.Auth())
	rg.POST("", meetingCtrl.Create)
	rg.PUT("", meetingCtrl.Update)
	rg.GET("", meetingCtrl.Get)
	rg.DELETE("", meetingCtrl.Delete)
	rg.GET("/list", meetingCtrl.List)

	rg.POST("/upload_resume", meetingCtrl.UploadResume)
	rg.POST("/ai_interview", meetingCtrl.AIInterview)
	rg.GET("/remark", meetingCtrl.GetRemark)
}
