package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"varchar(11);not null;unique"`
	Password  string `gorm:"size:255;not null"`
}

func main() {
	db := InitDB()
	r := gin.Default()
	r.POST("/api/auth/register", func(c *gin.Context) {
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
			name = RandomString(10)
		}

		log.Println(name, password, telephone) // 验证注册内容输出
		// 判断手机号是否存在
		if isTelephoneExist(db, telephone) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户已存在"})
			return
		}

		// 创建用户
		newUser := User{
			Name:      name,
			Telephone: telephone,
			Password:  password,
		}
		db.Create(&newUser)

		// 返回结果
		c.JSON(200, gin.H{
			"msg": "注册成功",
		})
	})
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

// 生成随机数
func RandomString(n int) string {
	var letters = []byte("asdfghjklqwertyuiopQWERTYUIOPASDFGHJKL1234567890")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// 连接数据库
func InitDB() *gorm.DB {
	host := "localhost"
	port := "3306"
	database := "govue"
	username := "root"
	password := "wangshaojun"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)
	db, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		panic("failed tp connect to database, err: " + err.Error())
	}
	db.AutoMigrate(&User{}) // 自动创建表

	return db
}

// 判断手机号是否存在
func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
