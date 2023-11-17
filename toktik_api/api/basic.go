package api

import (
	"mirco_tiktok/toktik_api/forms"
	proto "mirco_tiktok/toktik_api/proto"
)

func ProtoUserInfoToFormsUserInfo(userinfo *proto.UserInfo) *forms.UserInfo {
	return &forms.UserInfo{
		Id:              userinfo.Id,
		Name:            userinfo.Name,
		FollowCount:     userinfo.FollowCount,
		FollowerCount:   userinfo.FollowerCount,
		IsFollow:        userinfo.IsFollow,
		Avatar:          userinfo.Avatar,
		BackgroundImage: userinfo.BackgroundImage,
		Signature:       userinfo.Signature,
		TotalFavorited:  userinfo.TotalFavorited,
		WorkCount:       userinfo.WorkCount,
		FavoriteCount:   userinfo.FavoriteCount,
	}
}

func ProtoVideoInfoToFormsVideo(videoInfo *proto.VideoInfo) *forms.Video {
	author := ProtoUserInfoToFormsUserInfo(videoInfo.Author)
	return &forms.Video{
		Id:            videoInfo.Id,
		UserInfoId:    videoInfo.User_InfoId,
		Author:        *author,
		PlayUrl:       videoInfo.PlayUrl,
		CoverUrl:      videoInfo.PlayUrl,
		FavoriteCount: videoInfo.FavoriteCount,
		CommentCount:  videoInfo.CommentCount,
		IsFavorite:    videoInfo.IsFavorite,
		Title:         videoInfo.Title,
	}
}

func ProtoCommentInfoToFormsCommentInfo(commentInfo *proto.CommentInfo) *forms.Comment {
	user := ProtoUserInfoToFormsUserInfo(commentInfo.UserInfo)
	return &forms.Comment{
		Id:         commentInfo.Id,
		UserInfoId: commentInfo.UserInfoId,
		VideoId:    commentInfo.VideoId,
		User:       *user,
		Content:    commentInfo.Content,
		CreateDate: commentInfo.CreateDate,
	}
}
