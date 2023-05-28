package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"uisee.com/govue/common"
	"uisee.com/govue/model"
)

func AuthMiddleeware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取 authorization header
		tokenString := ctx.GetHeader("Authorization")

		// 验证格式
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") { // token为空或者不是Bearer开头说明token错误
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort() //抛弃请求
			return
		}

		tokenString = tokenString[7:] //提取token有效部分,截取掉Bearer加空格

		// 解析token
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid { // 解析失败
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort() //抛弃请求
			return
		}

		// token 通过验证，获取到了token中的userID
		userId := claims.UserId
		DB := common.GetDB()
		var user model.User
		DB.First(&user, userId)

		// 验证用户是否存在
		if userId == 0 { // 用户不存在 token 无效
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort() //抛弃请求
			return
		}

		// 用户存在，将user信息，写入上下文
		ctx.Set("user", user)
		ctx.Next()

	}
}
