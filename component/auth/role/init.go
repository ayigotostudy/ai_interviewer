package role

import (
	"ai_jianli_go/config"
	"ai_jianli_go/logs"

	"github.com/casbin/casbin/v2"
)

// 定义角色关系
type Role int

// 用户角色
const (
	Guest       = iota // 游客
	Common             // 普通用户
	Member             // 会员
	SuperMember        // 超级会员
	SuperAdmin         // 超级管理员
)

var roleMap = map[Role]string{
	Guest:       "guest",
	Common:      "common",
	Member:      "member",
	SuperMember: "super_member",
	SuperAdmin:  "super_admin",
}

var enforcer *casbin.Enforcer

func InitCasbin() {
	var err error
	enforcer, err = casbin.NewEnforcer(config.GetRoleConfig().Model, config.GetRoleConfig().Policy)
	if err != nil {
		logs.SugarLogger.Error("初始化casbin错误，errs：" + err.Error())
	}
}

// 用户身份int转对应的string
func GetRoleString(r int64) string {
	if role, ok := roleMap[Role(r)]; ok {
		return role
	}
	return roleMap[Guest]
}

// 获取所有角色列表
func GetAllRoles() map[Role]string {
	return roleMap
}

// 根据角色名称获取角色ID
func GetRoleByName(roleName string) Role {
	for role, name := range roleMap {
		if name == roleName {
			return role
		}
	}
	return Guest // 默认返回游客角色
}

// 检查角色是否存在
func IsValidRole(role Role) bool {
	_, exists := roleMap[role]
	return exists
}

// 获取角色描述
func GetRoleDescription(role Role) string {
	descriptions := map[Role]string{
		Guest:       "游客 - 只能访问公开接口和登录注册",
		Common:      "普通用户 - 可以管理个人资料、简历、会议和知识库",
		Member:      "会员 - 可以使用AI聊天、分析和语音识别功能",
		SuperMember: "超级会员 - 可以使用高级AI功能和批量处理",
		SuperAdmin:  "超级管理员 - 拥有系统管理权限",
	}
	return descriptions[role]
}

// 获取Enforcer实例
func GetEnforcer() *casbin.Enforcer {
	return enforcer
}
