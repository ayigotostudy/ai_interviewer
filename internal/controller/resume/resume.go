package resumeController

import (
	"ai_jianli_go/internal/controller"
	resumeService "ai_jianli_go/internal/service/resume"
	"ai_jianli_go/types/req"
	"ai_jianli_go/types/resp/common"
	"context"
	"strconv"

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

func (mc *ResumeController) GetResumeTemplate(c *gin.Context) {
	ctrl := controller.NewCtrl[req.NoReq](c)
	if err := c.Bind(ctrl.Request); err != nil {
		ctrl.NoDataJSON(common.CodeInvalidParams)
		return
	}
	resume, code := mc.svc.GetResumeTemplate(context.Background())
	ctrl.WithDataJSON(code, resume)
}

func (mc *ResumeController) GetResumeList(c *gin.Context) {
	ctrl := controller.NewCtrl[req.GetResumeListRequest](c)
	ctrl.Request.UserID = c.GetUint("id")
	resumes, code := mc.svc.GetResumeList(context.Background(), ctrl.Request.UserID)
	ctrl.WithDataJSON(code, resumes)
}

func (mc *ResumeController) GetResume(c *gin.Context) {
	ctrl := controller.NewCtrl[req.GetResumeDetailRequest](c)
	idStr := c.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctrl.NoDataJSON(common.CodeInvalidParams)
		return
	}
	ctrl.Request.ID = uint(id)

	resume, code := mc.svc.GetResume(context.Background(), ctrl.Request.ID)
	ctrl.WithDataJSON(code, resume)
}


func (mc *ResumeController) DeleteResume(c *gin.Context) {
	ctrl := controller.NewCtrl[req.DeleteResumeRequest](c)
	if err := c.Bind(ctrl.Request); err != nil {
		ctrl.NoDataJSON(common.CodeInvalidParams)
		return
	}
	code := mc.svc.DeleteResume(context.Background(), ctrl.Request.ID)
	ctrl.WithDataJSON(code, nil)
}

func (mc *ResumeController) UpdateResume(c *gin.Context) {
	ctrl := controller.NewCtrl[req.UpdateResumeRequest](c)
	if err := c.Bind(ctrl.Request); err != nil {
		ctrl.NoDataJSON(common.CodeInvalidParams)
		return
	}
	code := mc.svc.UpdateResume(context.Background(), ctrl.Request)
	ctrl.WithDataJSON(code, nil)
}