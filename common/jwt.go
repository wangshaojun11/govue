package common

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"uisee.com/govue/model"
)

// 定义 jwt 密钥
var jwtKey = []byte("a_secret_crect")

// 定义token的claims
type Claims struct {
	UserId uint
	jwt.StandardClaims
}

// 登录成功调用方法发放token
func ReleaseToken(user model.User) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour) // 定义过期时间
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // token过期时间
			IssuedAt:  time.Now().Unix(),     // token 发放时间
			Issuer:    "oceanlearn.tech",     // token 发放人
			Subject:   "user.token",          // 主题

		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey) //使用jwt密钥生成token
	// token 生成错误返回错误
	if err != nil {
		return "", err 
	}
	return tokenString, nil
}


// 解析token函数, 从tokenString 中解析出 claims
func ParseToken(tokenString string) (*jwt.Token, *Claims, error)  {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, claims, err
}

