package router

import (
	"ai_jianli_go/component"
	wikiController "ai_jianli_go/internal/controller/wiki"
	"ai_jianli_go/internal/dao"
	"ai_jianli_go/internal/middleware"
	wikiService "ai_jianli_go/internal/service/wiki"

	"github.com/gin-gonic/gin"
)

func wiki(r *gin.RouterGroup) {
	ctrl := wikiController.NewWikiController(wikiService.NewWikiService(dao.NewWikiDAO(component.GetMySQLDB())))
	r.Use(middleware.Auth())
	r.POST("/", ctrl.CreateWiki)
	r.GET("/list", ctrl.GetWikiList)
	r.GET("/", ctrl.GetWiki)
	r.DELETE("/", ctrl.DeleteWiki)
	r.POST("/query", ctrl.QueryWiki)
}