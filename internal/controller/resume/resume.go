package resumeController

import (
	"ai_jianli_go/internal/controller"
	resumeService "ai_jianli_go/internal/service/resume"
	"ai_jianli_go/types/req"
	"ai_jianli_go/types/resp/common"
	"context"

	"github.com/gin-gonic/gin"
)

type ResumeController struct {
	svc *resumeService.ResumeService
}

func NewResumeController(svc *resumeService.ResumeService) *ResumeController {
	return &ResumeController{
		svc: svc,
	}
}

func (mc *ResumeController) CreateResume(c *gin.Context) {
	ctrl := controller.NewCtrl[req.CreateResumeRequest](c)
	if err := c.Bind(ctrl.Request); err != nil {
		ctrl.NoDataJSON(common.CodeInvalidParams)
		return
	}
	ctrl.Request.UserID = c.GetUint("id")
	resume, code := mc.svc.CreateResume(context.Background(), ctrl.Request)
	ctrl.WithDataJSON(code, resume)
}
