package router

import (
	"ai_jianli_go/component"
	userController "ai_jianli_go/internal/controller/user"
	"ai_jianli_go/internal/dao"
	userService "ai_jianli_go/internal/service/user"

	"github.com/gin-gonic/gin"
)

func user(r *gin.RouterGroup) {
	ctrl := userController.NewUserController(userService.NewUserService(dao.NewUserDAO(component.GetMySQLDB())))
	r.POST("/login", ctrl.Login)
	r.POST("/register", ctrl.Register)
}
