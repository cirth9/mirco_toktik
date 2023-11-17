package models

import (
	"gorm.io/gorm"
	"mirco_tiktok/toktik_srv/global"
	"strconv"
)

type UserFavorVideos struct {
	UserInfoId int64
	VideoId    int64
}

func AddUserFavorVideos(userInfoId, VideoId int64) {
	userFavorVideos := UserFavorVideos{
		UserInfoId: userInfoId,
		VideoId:    VideoId,
	}
	global.DB.Create(&userFavorVideos)
}

func RmUserFavorVideos(userInfoId, videoId int64) {
	global.DB.Where("user_info_id = ? and video_id = ?", userInfoId, videoId).Delete(&UserFavorVideos{})
}

func IsUserFavorVideo(userInfoId, VideoId int64) bool {
	if global.DB.
		Where("user_info_id = ? and video_id = ?", userInfoId, VideoId).
		First(&UserFavorVideos{}).Error == gorm.ErrRecordNotFound {
		return false
	}
	return true
}

func GetUserFavorVideoId(userId int64) []string {
	var userFavorVideosList []*UserFavorVideos
	var favorVideoIdList []string
	global.DB.Where("user_info_id = ? ", userId).Find(&userFavorVideosList)
	for _, v := range userFavorVideosList {
		videoId := strconv.FormatInt(v.VideoId, 10)
		favorVideoIdList = append(favorVideoIdList, videoId)
	}
	return favorVideoIdList
}

func GetUserFavorList(userID int64) []*Video {
	var FavorVideoList []*Video
	VideoIdList := GetUserFavorVideoId(userID)
	global.DB.Where("id in ?", VideoIdList).Find(&FavorVideoList)
	return FavorVideoList
}
