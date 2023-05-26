package common

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"uisee.com/govue/model"
)

var DB *gorm.DB

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
	db.AutoMigrate(&model.User{}) // 自动创建表

	DB = db
	return db
}

// 获取DB实例
func GetDB() *gorm.DB {
	return DB
}
