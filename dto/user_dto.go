package dto

import "uisee.com/govue/model"

// 用户返回的结构体
type UserDto struct {
	Name string `json:"name"`
	Telephone string `json:"telephone"`
}

// 转换函数
func ToUserDto(user model.User) UserDto {
	return UserDto{
		Name: user.Name,
		Telephone: user.Telephone,
	}
}