package userController

import (
	"ai_jianli_go/internal/controller"
	userService "ai_jianli_go/internal/service/user"
	"ai_jianli_go/types/req"
	"ai_jianli_go/types/resp/common"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	svc *userService.UserService
}

func NewUserController(s *userService.UserService) *UserController {
	return &UserController{
		svc: s,
	}
}

func (u *UserController) Register(c *gin.Context) {
	ctrl := controller.NewCtrl[req.RegisterReq](c)
	if err := c.Bind(ctrl.Request); err != nil {
		ctrl.NoDataJSON(common.CodeInvalidParams)
		return
	}

	code := u.svc.Register(ctrl.Request)
	ctrl.NoDataJSON(code)
}

func (u *UserController) Login(c *gin.Context) {
	ctrl := controller.NewCtrl[req.LoginReq](c)
	if err := c.Bind(ctrl.Request); err != nil {
		ctrl.NoDataJSON(common.CodeInvalidParams)
		return
	}

	data, code := u.svc.Login(ctrl.Request)
	ctrl.WithDataJSON(code, data)
}
