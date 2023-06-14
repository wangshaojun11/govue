package main

import (
	"github.com/gin-gonic/gin"
	"uisee.com/govue/controller"
	"uisee.com/govue/middleware"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware())                                     // 跨域
	r.POST("/api/auth/register", controller.Register)                      // 注册
	r.POST("/api/auth/login", controller.Login)                            // 登录
	r.GET("/api/auth/info", middleware.AuthMiddleeware(), controller.Info) // 用户信息,中间件保护用户信息接口
	return r
}
