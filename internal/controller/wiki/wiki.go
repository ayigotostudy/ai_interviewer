package wikiController

import (
	"ai_jianli_go/config"
	"ai_jianli_go/internal/controller"
	wikiService "ai_jianli_go/internal/service/wiki"
	"ai_jianli_go/logs"
	"ai_jianli_go/types/model"
	"ai_jianli_go/types/req"
	"ai_jianli_go/types/resp/common"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type WikiController struct {
	svc *wikiService.WikiService
}

func NewWikiController(svc *wikiService.WikiService) *WikiController {
	return &WikiController{svc: svc}
}

func (c *WikiController) CreateWiki(ctx *gin.Context) {
	ctrl := controller.NewCtrl[req.CreateWikiRequest](ctx)

	// 直接赋值form数据
	ctrl.Request.Title = ctx.PostForm("title")
	ctrl.Request.Content = ctx.PostForm("content")
	ctrl.Request.WikiType, _ = strconv.Atoi(ctx.PostForm("wiki_type"))
	ctrl.Request.Type, _ = strconv.Atoi(ctx.PostForm("type"))
	rootId, _ := strconv.ParseUint(ctx.PostForm("root_id"), 10, 64)
	ctrl.Request.RootId = uint(rootId)
	ctrl.Request.Url = ctx.PostForm("url")
	ctrl.Request.UserID = ctx.GetUint("id")
	parentId, _ := strconv.ParseUint(ctx.PostForm("parent_id"), 10, 64)
	ctrl.Request.ParentID = uint(parentId)

	if ctrl.Request.Type == model.WikiTypeArticle {
		// 检查是否有直接传递的URL
		if ctrl.Request.Url != "" {
			logs.SugarLogger.Infof("使用直接传递的URL: %s", ctrl.Request.Url)
		} else {
			// 处理文件上传
			file, err := ctx.FormFile("file")
			if err != nil {
				logs.SugarLogger.Errorf("文件上传失败: %v", err)
				ctrl.NoDataJSON(common.CodeInvalidParams)
				return
			}

			// 保存上传的文件
			lname := path.Ext(file.Filename)
			file.Filename = ctrl.Request.Title + "_" + time.Now().Format("20060102150405") + lname
			workdir, _ := os.Getwd()
			savepath := filepath.Join(workdir, config.GetLocalPathConfig().Path, "wiki", file.Filename)
			err = ctx.SaveUploadedFile(file, savepath)
			if err != nil {
				logs.SugarLogger.Errorf("保存文件失败: %v", err)
				ctrl.NoDataJSON(common.CodeInvalidParams)
				return
			}
			ctrl.Request.Url = savepath
			logs.SugarLogger.Infof("文件上传成功，保存路径: %s", savepath)
		}
	}

	code := c.svc.CreateWiki(ctrl.Request)
	ctrl.NoDataJSON(code)
}

func (c *WikiController) GetWikiList(ctx *gin.Context) {
	ctrl := controller.NewCtrl[req.GetWikiListRequest](ctx)
	ctrl.Request.UserID = ctx.GetUint("id")
	wikis, code := c.svc.GetWikiList(ctrl.Request)
	ctrl.WithDataJSON(code, wikis)
}

func (c *WikiController) GetWiki(ctx *gin.Context) {
	ctrl := controller.NewCtrl[req.GetWikiRequest](ctx)
	id, _ := strconv.ParseUint(ctx.Query("id"), 10, 64)
	ctrl.Request.ID = uint(id)
	ctrl.Request.UserID = ctx.GetUint("id")
	wiki, code := c.svc.GetWiki(ctrl.Request)
	ctrl.WithDataJSON(code, wiki)
}

func (c *WikiController) DeleteWiki(ctx *gin.Context) {
	ctrl := controller.NewCtrl[req.DeleteWikiRequest](ctx)
	if err := ctx.Bind(ctrl.Request); err != nil {
		ctrl.NoDataJSON(common.CodeInvalidParams)
		return
	}
	ctrl.Request.UserID = ctx.GetUint("id")
	code := c.svc.DeleteWiki(ctrl.Request)
	ctrl.NoDataJSON(code)
}

func (c *WikiController) QueryWiki(ctx *gin.Context) {
	ctrl := controller.NewCtrl[req.QueryWikiRequest](ctx)
	if err := ctx.Bind(ctrl.Request); err != nil {
		ctrl.NoDataJSON(common.CodeInvalidParams)
		return
	}
	ctrl.Request.UserID = ctx.GetUint("id")
	data, code := c.svc.Query(ctrl.Request)
	ctrl.WithDataJSON(code, data)
}

func (c *WikiController) GetListByParentId(ctx *gin.Context) {
	ctrl := controller.NewCtrl[req.GetWikiListRequest](ctx)
	ctrl.Request.UserID = ctx.GetUint("id")
	parentId, _ := strconv.ParseUint(ctx.Query("parent_id"), 10, 64)
	ctrl.Request.ParentID = uint(parentId)
	wikis, code := c.svc.GetListByParentId(ctrl.Request)
	ctrl.WithDataJSON(code, wikis)
}

func (c *WikiController) GetFileByPath(ctx *gin.Context) {
	// 获取文件路径参数
	filePath := ctx.Query("path")
	if filePath == "" {
		// 返回参数错误
		controller.NewCtrl[interface{}](ctx).NoDataJSON(common.CodeInvalidParams)
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 文件不存在，返回错误
		controller.NewCtrl[interface{}](ctx).NoDataJSON(common.CodeFileNotFound)
		return
	} else if err != nil {
		// 其他错误
		logs.SugarLogger.Errorf("检查文件状态失败: %v", err)
		controller.NewCtrl[interface{}](ctx).NoDataJSON(common.CodeFileNotFound)
		return
	}

	// 获取文件名用于设置Content-Disposition
	fileName := filepath.Base(filePath)

	// 设置响应头，确保中文文件名正确显示
	disposition := fmt.Sprintf("attachment; filename=\"%s\"; filename*=UTF-8''%s",
		fileName, url.QueryEscape(fileName))
	ctx.Header("Content-Disposition", disposition)

	// 返回文件内容
	ctx.File(filePath)
}