package meetingController

import (
	"ai_jianli_go/internal/controller"
	meetingService "ai_jianli_go/internal/service/meeting"
	"ai_jianli_go/types/req"
	"ai_jianli_go/types/resp/common"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MeetingController struct {
	svc *meetingService.MeetingService
}

func NewMeetingController(s *meetingService.MeetingService) *MeetingController {
	return &MeetingController{svc: s}
}

func (mc *MeetingController) Create(c *gin.Context) {
	ctrl := controller.NewCtrl[req.CreateMeetingReq](c)
	if err := c.Bind(ctrl.Request); err != nil {
		ctrl.NoDataJSON(common.CodeInvalidParams)
		return
	}
	ctrl.Request.UserID = c.GetUint("id")
	code := mc.svc.Create(ctrl.Request)
	ctrl.NoDataJSON(code)
}

func (mc *MeetingController) Update(c *gin.Context) {
	ctrl := controller.NewCtrl[req.UpdateMeetingReq](c)
	if err := c.Bind(ctrl.Request); err != nil {
		ctrl.NoDataJSON(common.CodeInvalidParams)
		return
	}
	ctrl.Request.UserID = c.GetUint("id")
	code := mc.svc.Update(ctrl.Request)
	ctrl.NoDataJSON(code)
}

func (mc *MeetingController) Get(c *gin.Context) {
	ctrl := controller.NewCtrl[req.GetMeetingReq](c)

	id, err := strconv.ParseUint(c.Query("id"), 10, 64) // 参数：字符串, 进制(10), 位数(64)
	if err != nil {
		ctrl.NoDataJSON(common.CodeInvalidParams)
		return
	}

	ctrl.Request.ID = uint(id)

	meeting, code := mc.svc.Get(ctrl.Request.ID)
	if code != common.CodeSuccess {
		ctrl.NoDataJSON(code)
		return
	}
	ctrl.WithDataJSON(code, meeting)
}

func (mc *MeetingController) List(c *gin.Context) {
	meetings, code := mc.svc.List()
	if code != common.CodeSuccess {
		ctrl := controller.NewCtrl[any](c)
		ctrl.NoDataJSON(code)
		return
	}
	ctrl := controller.NewCtrl[any](c)
	ctrl.WithDataJSON(code, meetings)
}

func (mc *MeetingController) Delete(c *gin.Context) {
	ctrl := controller.NewCtrl[req.GetMeetingReq](c)
	if err := c.Bind(ctrl.Request); err != nil {
		ctrl.NoDataJSON(common.CodeInvalidParams)
		return
	}
	code := mc.svc.Delete(ctrl.Request.ID)
	ctrl.NoDataJSON(code)
}

// 上传简历接口
func (mc *MeetingController) UploadResume(c *gin.Context) {
	ctrl := controller.NewCtrl[req.UploadResumeReq](c)
	if err := c.Bind(ctrl.Request); err != nil {
		ctrl.NoDataJSON(common.CodeInvalidParams)
		return
	}
	ctrl.Request.UserID = c.GetUint("id")
	code := mc.svc.UploadResume(ctrl.Request)
	ctrl.NoDataJSON(code)
}

// AI面试接口
func (mc *MeetingController) AIInterview(c *gin.Context) {
	ctrl := controller.NewCtrl[req.AIInterviewReq](c)
	if err := c.Bind(ctrl.Request); err != nil {
		ctrl.NoDataJSON(common.CodeInvalidParams)
		return
	}
	ctrl.Request.UserID = c.GetUint("id")
	reply, code := mc.svc.AIInterview(ctrl.Request)
	if code != common.CodeSuccess {
		ctrl.NoDataJSON(code)
		return
	}
	ctrl.WithDataJSON(code, gin.H{"reply": reply})
}
