package main

import (
	"github.com/gin-gonic/gin"
	"uisee.com/govue/common"
)

func main() {
	common.InitDB()

	r := gin.Default()
	r = CollectRoute(r)
	panic(r.Run()) // 监听并在 0.0.0.0:8080 上启动服务
}
