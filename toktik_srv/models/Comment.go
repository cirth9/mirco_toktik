package models

import (
	"mirco_tiktok/toktik_srv/global"
	"time"
)

/*
message Comment {
  required int64 id = 1; // 视频评论id
  required User user =2; // 评论用户信息
  required string content = 3; // 评论内容
  required string create_date = 4; // 评论发布日期，格式 mm-dd
}
*/

type Comment struct {
	Id         int64     `json:"id"`
	UserInfoId int64     `json:"-"` //用于一对多关系的id
	VideoId    int64     `json:"-"` //一对多，视频对评论
	User       UserInfo  `json:"user" gorm:"-"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"-"`
	CreateDate string    `json:"create_date" gorm:"-"`
}

func (c *Comment) TableName() string {
	return "comment"
}

func AddNewComment(comment Comment) {
	global.DB.Create(&comment)
}

func RemoveComment(commentId int64) {
	global.DB.Where("id = ?", commentId).Delete(&Comment{})
}

func GetAllCommentByVid(videoId int64) []*Comment {
	var commentList []*Comment
	global.DB.Where("video_id = ?", videoId).Find(&commentList)
	return commentList
}
