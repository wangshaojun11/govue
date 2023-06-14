package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 通过这个中间给header 写入允许访问的域名
func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080") // 前台域名
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")                      // 缓存时间
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "*")                    // 方法
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")                    // Header
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")             // 凭证

		// 判断请求方法是否是 options
		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(200)
		}else {
			ctx.Next() // 不是 Option 向下传递
		}

	}
}
