package main

import (
	"ai_jianli_go/component"
	"ai_jianli_go/component/auth/role"
	"ai_jianli_go/config"
	"ai_jianli_go/internal/router"
	"ai_jianli_go/logs"
	"ai_jianli_go/pkg/rag"
)

func main() {
	logs.Init()
	config.Init()
	component.Init()
	rag.Init()
	role.InitCasbin()

	router := router.Init()

	router.Run(":8080")
}
