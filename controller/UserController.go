package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"uisee.com/govue/common"
	"uisee.com/govue/dto"
	"uisee.com/govue/model"
	"uisee.com/govue/response"
	"uisee.com/govue/util"
)

func Register(ctx *gin.Context) {
	DB := common.GetDB()

	// 方法1: 使用map
	// var requestMap = make(map[string]string)
	// json.NewDecoder(ctx.Request.Body).Decode(&requestMap)

	// 方法2:使用结构体承接参数
	// var requestUser = model.User{}
	// json.NewDecoder(ctx.Request.Body).Decode(&requestUser)

	// 方法3:使用bind方法
	var requestUser = model.User{}
	ctx.Bind(&requestUser)

	// 获取参数  创建的时候传3个参数，姓名，手机号，密码
	name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password

	// 数据验证
	if len(telephone) != 11 { // 验证手机号是否为11位
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 { // 验证密码是否少于6位
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码大于6位")
		return
	}
	if len(name) == 0 { // 验证用户为空，生成随机10位字符
		name = util.RandomString(10)
	}

	log.Println(name, password, telephone) // 验证注册内容输出
	// 判断手机号是否存在
	if isTelephoneExist(DB, telephone) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户已存在")
		return
	}

	// 创建用户
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // 密码加密
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 500, nil, "用户密码加密失败")
		return
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}
	DB.Create(&newUser)
	// 发放token
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统错误")
		log.Printf("token generate err: %v", err)
		return
	}

	// 返回结果
	response.Success(ctx, gin.H{"token": token}, "注册成功")

}

// 用户登录
func Login(ctx *gin.Context) {
	DB := common.GetDB()
	// 获取参数
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	// 数据验证
	if len(telephone) != 11 { // 验证手机号是否为11位
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 { // 验证密码是否少于6位
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码大于6位")
		return
	}

	// 判断手机号是否存在
	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}

	// 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "密码错误")
		return
	}

	// 发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统错误")
		log.Printf("token generate err: %v", err)
		return
	}

	// 返回结果
	response.Success(ctx, gin.H{"token": token}, "登录成功")
}

// 获取用户信息
func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"user": dto.ToUserDto(user.(model.User)),
		},
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
