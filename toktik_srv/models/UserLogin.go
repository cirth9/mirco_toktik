package models

import (
	"mirco_tiktok/toktik_srv/global"
)

type UserLogin struct {
	Id         int64 `gorm:"primary_key"`
	UserInfoId int64
	Username   string `gorm:"primary_key"`
	Password   string `gorm:"password"`
}

func (user *UserLogin) TableName() string {
	return "user_login"
}

func FindUsername(username string) UserLogin {
	userLogin := UserLogin{}
	global.DB.Where("username = ?", username).First(&userLogin)
	return userLogin
}

func CreateUserLogin(userLogin UserLogin) int64 {
	global.DB.Create(&userLogin)
	return userLogin.Id
}
