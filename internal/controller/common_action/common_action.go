package common_action

import (
	"ai_jianli_go/internal/controller"
	"ai_jianli_go/internal/service/common_action"

	"ai_jianli_go/types/req"

	"github.com/gin-gonic/gin"
)

type CommonActionController struct {
	svc *common_action.CommonActionService
}

func NewCommonActionController(svc *common_action.CommonActionService) *CommonActionController {
	return &CommonActionController{svc: svc}
}

func (c *CommonActionController) SendAuthCode(ctx *gin.Context) {
	ctrl := controller.NewCtrl[req.NoReq](ctx)
	ctrl.NoDataJSON(c.svc.SendAuthCode(ctx.Query("email")))
}