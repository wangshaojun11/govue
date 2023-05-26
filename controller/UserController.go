package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"uisee.com/govue/common"
	"uisee.com/govue/model"
	"uisee.com/govue/util"
)

func Register(c *gin.Context) {
	DB := common.GetDB()
	// 获取参数
	name := c.PostForm("name") // 创建的时候传3个参数，姓名，手机号，密码
	telephone := c.PostForm("telephone")
	password := c.PostForm("password")

	// 数据验证
	if len(telephone) != 11 { // 验证手机号是否为11位
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
		return
	}
	if len(password) < 6 { // 验证密码是否少于6位
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码大于6位"})
		return
	}
	if len(name) == 0 { // 验证用户为空，生成随机10位字符
		name = util.RandomString(10)
	}

	log.Println(name, password, telephone) // 验证注册内容输出
	// 判断手机号是否存在
	if isTelephoneExist(DB, telephone) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户已存在"})
		return
	}

	// 创建用户
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  password,
	}
	DB.Create(&newUser)

	// 返回结果
	c.JSON(200, gin.H{
		"msg": "注册成功",
	})
}

// 判断手机号是否存在
func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
