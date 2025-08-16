package router

import (
	"ai_jianli_go/component"
	"ai_jianli_go/internal/dao"
    "ai_jianli_go/internal/service/resume"
	"ai_jianli_go/internal/controller/resume"
	"github.com/gin-gonic/gin"
)

func resume(r *gin.RouterGroup) {
	resumeDAO := dao.NewResumeDAO(component.GetMySQLDB())
	resumeSvc := resumeService.NewResumeService(resumeDAO)
	resumeCtrl := resumeController.NewResumeController(resumeSvc)
	r.POST("", resumeCtrl.CreateResume)
}
