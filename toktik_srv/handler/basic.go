package handler

import (
	models2 "mirco_tiktok/toktik_srv/models"
	proto "mirco_tiktok/toktik_srv/proto"
)

func VideoToProtoVideoInfo(video *models2.Video) *proto.VideoInfo {
	author := models2.FindUserInfoById(video.UserInfoId)
	authorInfo := &proto.UserInfo{
		Id:              author.Id,
		Name:            author.Name,
		FollowCount:     author.FollowCount,
		FollowerCount:   author.FollowerCount,
		IsFollow:        author.IsFollow,
		Avatar:          author.Avatar,
		BackgroundImage: author.BackgroundImage,
		Signature:       author.Signature,
		TotalFavorited:  author.TotalFavorited,
		WorkCount:       author.WorkCount,
		FavoriteCount:   author.FavoriteCount,
	}
	return &proto.VideoInfo{
		Id:            video.Id,
		User_InfoId:   video.UserInfoId,
		Author:        authorInfo,
		PlayUrl:       video.PlayUrl,
		CoverUrl:      video.CoverUrl,
		FavoriteCount: video.FavoriteCount,
		CommentCount:  video.CommentCount,
		IsFavorite:    video.IsFavorite,
		Title:         video.Title,
	}
}

func UserToProtoUserInfo(user *models2.UserInfo) *proto.UserInfo {
	return &proto.UserInfo{
		Id:              user.Id,
		Name:            user.Name,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		IsFollow:        user.IsFollow,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}
}

func CommentToProtoCommentInfo(comment *models2.Comment) *proto.CommentInfo {
	userInfo := UserToProtoUserInfo(&comment.User)
	return &proto.CommentInfo{
		Id:         comment.Id,
		UserInfoId: comment.UserInfoId,
		VideoId:    comment.VideoId,
		Content:    comment.Content,
		CreateDate: comment.CreateDate,
		UserInfo:   userInfo,
	}
}
