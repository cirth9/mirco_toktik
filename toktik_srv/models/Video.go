package models

import (
	"gorm.io/gorm"
	"log"
	"mirco_tiktok/toktik_srv/global"
	"time"
)

/*
	message Video {
	  required int64 id = 1; // 视频唯一标识
	  required User author = 2; // 视频作者信息
	  required string play_url = 3; // 视频播放地址
	  required string cover_url = 4; // 视频封面地址
	  required int64 favorite_count = 5; // 视频的点赞总数
	  required int64 comment_count = 6; // 视频的评论总数
	  required bool is_favorite = 7; // true-已点赞，false-未点赞
	  required string title = 8; // 视频标题
	}
*/

var (
	maxLimit = 30
)

type Video struct {
	Id            int64       `json:"id,omitempty"`
	UserInfoId    int64       `json:"-"`
	Author        UserInfo    `json:"author,omitempty" gorm:"-"`
	PlayUrl       string      `json:"play_url,omitempty"`
	CoverUrl      string      `json:"cover_url,omitempty"`
	FavoriteCount int64       `json:"favorite_count,omitempty"`
	CommentCount  int64       `json:"comment_count,omitempty"`
	IsFavorite    bool        `json:"is_favorite,omitempty"`
	Title         string      `json:"title,omitempty"`
	Users         []*UserInfo `json:"-" gorm:"many2many:user_favor_videos;"`
	Comments      []*Comment  `json:"-"`
	CreatedAt     time.Time   `json:"-"`
	UpdatedAt     time.Time   `json:"-"`
}

func (v *Video) TableName() string {
	return "video"
}

func CreateVideo(video Video) *gorm.DB {
	return global.DB.Create(&video)
}

func FindVideoByUid(userId int64) []*Video {
	var videoList []*Video
	global.DB.Where("user_info_id = ?", userId).Find(&videoList)
	return videoList
}

func GetVideoList(latestTime time.Time) []*Video {
	var videoFeed []*Video
	global.DB.
		Where("created_at < ?", latestTime).
		Order("created_at asc").
		Limit(maxLimit).
		Find(&videoFeed)
	for _, v := range videoFeed {
		author := FindUserInfoById(v.UserInfoId)
		v.Author = author
	}
	return videoFeed
}

func GetVideoListWithoutSelf(latestTime time.Time, userId int64) []*Video {
	var videoFeed []*Video
	global.DB.
		Where("created_at < ? and user_info_id != ?", latestTime, userId).
		Order("created_at asc").
		Limit(maxLimit).
		Find(&videoFeed)
	log.Printf("%+v", videoFeed)
	for _, v := range videoFeed {
		author := FindUserInfoById(v.UserInfoId)
		v.Author = author
	}
	return videoFeed
}

func AscFavoriteCountById(videoID int64) {
	global.DB.Model(Video{}).
		Where("id = ?", videoID).
		Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
}

func DscFavoriteCountById(videoID int64) {
	global.DB.Model(Video{}).
		Where("id = ?", videoID).
		Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
}

func AscCommentCountById(videoID int64) {
	global.DB.Model(Video{}).
		Where("id = ?", videoID).
		Update("comment_count", gorm.Expr("comment_count + ?", 1))
}

func DscCommentCountById(videoID int64) {
	global.DB.Model(Video{}).
		Where("id = ?", videoID).
		Update("comment_count", gorm.Expr("comment_count - ?", 1))
}
