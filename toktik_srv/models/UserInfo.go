package models

import (
	"log"
	"mirco_tiktok/toktik_srv/global"
)

/*
message User {
  required int64 id = 1; // 用户id
  required string name = 2; // 用户名称
  optional int64 follow_count = 3; // 关注总数
  optional int64 follower_count = 4; // 粉丝总数
  required bool is_follow = 5; // true-已关注，false-未关注
  optional string avatar = 6; //用户头像
  optional string background_image = 7; //用户个人页顶部大图
  optional string signature = 8; //个人简介
  optional int64 total_favorited = 9; //获赞数量
  optional int64 work_count = 10; //作品数量
  optional int64 favorite_count = 11; //点赞数量
}
*/

type UserInfo struct {
	Id              int64       `json:"id" gorm:"id,omitempty"`
	Name            string      `json:"name" gorm:"name,omitempty"`
	FollowCount     int64       `json:"follow_count" gorm:"follow_count,omitempty"`
	FollowerCount   int64       `json:"follower_count" gorm:"follower_count,omitempty"`
	IsFollow        bool        `json:"is_follow" gorm:"is_follow,omitempty"`
	Avatar          string      `json:"avatar"`
	BackgroundImage string      `json:"background_image"`
	Signature       string      `json:"signature"`
	TotalFavorited  int64       `json:"total_favorited"`
	WorkCount       int64       `json:"work_count"`
	FavoriteCount   int64       `json:"favorite_count"`
	User            *UserLogin  `json:"-"`                                     //用户与账号密码之间的一对一
	Videos          []*Video    `json:"-"`                                     //用户与投稿视频的一对多
	Follows         []*UserInfo `json:"-" gorm:"many2many:user_relations;"`    //用户之间的多对多
	FavorVideos     []*Video    `json:"-" gorm:"many2many:user_favor_videos;"` //用户与点赞视频之间的多对多
	Comments        []*Comment  `json:"-"`                                     //用户与评论的一对多
}

func (user *UserInfo) TableName() string {
	return "user_info"
}

func CreateUserInfo(login UserLogin) int64 {
	userInfo := UserInfo{
		Name:            login.Username,
		FollowCount:     0,
		FollowerCount:   0,
		IsFollow:        false,
		Avatar:          "",
		BackgroundImage: "",
		Signature:       "",
		TotalFavorited:  0,
		WorkCount:       0,
		FavoriteCount:   0,
	}
	global.DB.Create(&userInfo)
	return userInfo.Id
}

func FindUserInfoById(userId int64) UserInfo {
	userInfo := UserInfo{}
	global.DB.Where("id = ?", userId).First(&userInfo)
	log.Println(userInfo)
	return userInfo
}

func IncreaseWorkCountByUid(userId int64) {
	var userInfo UserInfo
	global.DB.Where("id = ?", userId).First(&userInfo)
	userInfo.WorkCount += 1
	global.DB.Save(&userInfo)
}

func DecreaseWorkCountByUid(userId int64) {
	var userInfo UserInfo
	global.DB.Where("id = ?", userId).First(&userInfo)
	userInfo.WorkCount -= 1
	global.DB.Save(&userInfo)
}

func IncFollowCount(userId int64) {
	var userInfo UserInfo
	global.DB.Where("id = ? ", userId).First(&userInfo)
	userInfo.FollowCount++
	global.DB.Save(&userInfo)
}

func DecFollowCount(userId int64) {
	var userInfo UserInfo
	global.DB.Where("id = ? ", userId).First(&userInfo)
	userInfo.FollowCount--
	global.DB.Save(&userInfo)
}

func IncFollowerCount(userId int64) {
	var userInfo UserInfo
	global.DB.Where("id = ? ", userId).First(&userInfo)
	userInfo.FollowerCount++
	global.DB.Save(&userInfo)
}

func DecFollowerCount(userId int64) {
	var userInfo UserInfo
	global.DB.Where("id = ? ", userId).First(&userInfo)
	userInfo.FollowerCount--
	global.DB.Save(&userInfo)
}
