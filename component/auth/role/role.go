package role

import (
	"ai_jianli_go/logs"
	"ai_jianli_go/types/resp/common"
	"context"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)

var checkLock sync.Mutex

// 验证权限 - 通用权限检查
func CheckPermission(c context.Context, ctx *gin.Context, userId int64, role int64) (StatusCode int64) {
	userRole := GetRoleString(role)
	ctx.Set("role", userRole)
	return check(userId, userRole, ctx.FullPath(), string(ctx.Request.Method))
}

// 目前策略模型
func check(userId int64, sub, obj, act string) (StatusCode int64) {
	checkLock.Lock()
	defer checkLock.Unlock()

	ok, _ := enforcer.Enforce(sub, obj, act) // sub主体 , obj对象 , act动作
	if ok {
		return common.CodeSuccess
	}
	logs.SugarLogger.Error(fmt.Sprintf("权限不足,用户ID：%d, 角色：%s, 路径：%s, 请求方法：%s", userId, sub, obj, act))
	return common.CodeInvalidRoleAdmin
}
