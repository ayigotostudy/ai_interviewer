package main

import (
	"ai_jianli_go/internal/router"
	"ai_jianli_go/pkg/component"
)

func main() {
	component.Init()
	router := router.Init()

	router.Run(":8080")
}
