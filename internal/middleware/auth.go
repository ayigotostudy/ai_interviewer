package middleware

import (
	"ai_jianli_go/component/auth/role"
	"ai_jianli_go/logs"
	"ai_jianli_go/pkg/utils"
	"ai_jianli_go/types/resp/common"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 验证用户是否登录的中间件
func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res := common.Response{}
		t := ctx.GetHeader("Authorization") //得到字串开头
		if t == "" || !strings.HasPrefix(t, "Bearer ") {
			logs.SugarLogger.Errorf("认证失败，无效的Authorization头: %s", t)
			ctx.JSON(http.StatusUnauthorized, "bearer解析失败")
			ctx.Abort()
			return
		}

		t = t[7:]                          //扔掉头部
		tk, c, r, e := utils.ParseToken(t) //c为claim结构体的实例
		if e != nil || !tk.Valid {
			logs.SugarLogger.Errorf("认证失败，token解析错误: %v", e)
			res.SetNoData(common.CodeInvalidToken)
			ctx.JSON(http.StatusUnauthorized, res)
			ctx.Abort() //中间件不通过
			return
		}
		//查找用户
		//存储用户信息
		ctx.Set("id", c)


		// 认证用户角色权限
		StatusCode := role.CheckPermission(context.Background(), ctx, int64(c), int64(r))
		if StatusCode != common.CodeSuccess {
			res.SetNoData(StatusCode)
			ctx.JSON(http.StatusUnauthorized, res)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
