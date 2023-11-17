package handler

import (
	"context"
	"github.com/go-errors/errors"
	"go.uber.org/zap"
	models2 "mirco_tiktok/toktik_srv/models"
	"mirco_tiktok/toktik_srv/proto"
	utils2 "mirco_tiktok/toktik_srv/utils"
)

/*
type TokTikClient interface {
	// user api
	GetUserInfo(ctx context.Context, in *UserInfoRequest, opts ...grpc.CallOption) (*UserInfoResponse, error)
	UserRegister(ctx context.Context, in *UserInfoRequest, opts ...grpc.CallOption) (*UserBasicResponse, error)
	UserLogin(ctx context.Context, in *UserInfoRequest, opts ...grpc.CallOption) (*UserBasicResponse, error)
	CreateUserInfo(ctx context.Context, in *UserInfoRequest, opts ...grpc.CallOption) (*UserBasicResponse, error)
	// publish api
	PublishAction(ctx context.Context, in *UserPublishRequest, opts ...grpc.CallOption) (*UserBasicResponse, error)
	UserPublishList(ctx context.Context, in *UserPublishListRequest, opts ...grpc.CallOption) (*UserPublishListResponse, error)
	// favorite api
	UserFavoriteAction(ctx context.Context, in *UserFavoriteRequest, opts ...grpc.CallOption) (*UserBasicResponse, error)
	UserFavoriteList(ctx context.Context, in *UserFavoriteListRequest, opts ...grpc.CallOption) (*UserFavoriteListResponse, error)
	// comment api
	UserCommentAction(ctx context.Context, in *UserCommentRequest, opts ...grpc.CallOption) (*UserBasicResponse, error)
	UserCommentList(ctx context.Context, in *UserCommentListRequest, opts ...grpc.CallOption) (*UserCommentListResponse, error)
	// relation api
	UserRelationAction(ctx context.Context, in *UserRelationRequest, opts ...grpc.CallOption) (*UserBasicResponse, error)
	UserRelationFollowList(ctx context.Context, in *UserRelationFollowListRequest, opts ...grpc.CallOption) (*UserRelationListResponse, error)
	UserRelationFollowerList(ctx context.Context, in *UserRelationFollowerListRequest, opts ...grpc.CallOption) (*UserRelationListResponse, error)
}
*/

type TokTikServer struct {
	__.UnimplementedTokTikServer
}

func (toktikServer *TokTikServer) GetUserInfo(ctx context.Context, in *__.UserInfoRequest) (*__.UserInfoResponse, error) {
	userInfo := models2.FindUserInfoById(in.UserId)
	return &__.UserInfoResponse{
		Id:              userInfo.Id,
		Name:            userInfo.Name,
		FollowCount:     userInfo.FollowerCount,
		FollowerCount:   userInfo.FollowerCount,
		IsFollow:        userInfo.IsFollow,
		Avatar:          userInfo.Avatar,
		BackgroundImage: userInfo.BackgroundImage,
		Signature:       userInfo.Signature,
		TotalFavorited:  userInfo.TotalFavorited,
		WorkCount:       userInfo.WorkCount,
		FavoriteCount:   userInfo.FavoriteCount,
		UserBasicResponse: &__.UserBasicResponse{
			StatusCode: 0,
			StatusMsg:  "获取用户信息成功",
			UserId:     int32(userInfo.Id),
			Token:      "",
		},
	}, nil
}

func (toktikServer *TokTikServer) UserRegister(ctx context.Context, in *__.UserInfoRequest) (*__.UserBasicResponse, error) {
	userLogin := models2.FindUsername(in.Username)
	if userLogin.Username != "" {
		zap.S().Info("username has exited!", in.Username)
		return &__.UserBasicResponse{
			StatusCode: -1,
			StatusMsg:  "注册失败，已存在用户",
			UserId:     -1,
			Token:      "",
		}, errors.Errorf("注册失败，已存在用户")
	} else {
		userLogin.Username = in.Username
		userLogin.Password = utils2.EncryptPassword(in.Password)
		userID := models2.CreateUserInfo(userLogin)
		userLogin.UserInfoId = userID
		models2.CreateUserLogin(userLogin)
		token, err := utils2.GetToken(userLogin.Username, userLogin.Password, userID)
		if err != nil {
			zap.S().Info("can't get token", err)
			return &__.UserBasicResponse{
				StatusCode: -1,
				StatusMsg:  "注册失败，无法获取token",
				UserId:     -1,
				Token:      "",
			}, errors.Errorf("注册失败，无法获取token,err : ", err.Error())
		}
		zap.S().Info("注册成功,UserInfo: ", userLogin)
		return &__.UserBasicResponse{
			StatusCode: 0,
			StatusMsg:  "注册成功",
			UserId:     int32(userID),
			Token:      token,
		}, nil
	}
}

func (toktikServer *TokTikServer) UserLogin(ctx context.Context, in *__.UserInfoRequest) (*__.UserBasicResponse, error) {
	userLogin := models2.FindUsername(in.Username)
	if userLogin.Username == "" {
		zap.S().Info("用户名或者密码错误")
		return &__.UserBasicResponse{
			StatusCode: -1,
			StatusMsg:  "用户名或者密码错误",
			UserId:     int32(in.UserId),
			Token:      "",
		}, errors.Errorf("用户名或者密码错误")
	} else {
		if !utils2.DecryptPassword(in.Password, userLogin.Password) {
			zap.S().Info("用户名或者密码错误")
			return &__.UserBasicResponse{
				StatusCode: -1,
				StatusMsg:  "用户名或者密码错误",
				UserId:     int32(in.UserId),
				Token:      "",
			}, errors.Errorf("用户名或者密码错误")
		}
		token, err := utils2.GetToken(in.Username, in.Password, userLogin.UserInfoId)

		if err != nil {
			zap.S().Error("登陆成功，但是token生成失败,err: ", err)
			return &__.UserBasicResponse{
				StatusCode: -1,
				StatusMsg:  "登陆成功，但是token生成失败",
				UserId:     int32(in.UserId),
				Token:      "",
			}, errors.Errorf("登陆成功，但是token生成失败,err: ", err.Error())
		}
		return &__.UserBasicResponse{
			StatusCode: 0,
			StatusMsg:  "登陆成功！！",
			UserId:     int32(in.UserId),
			Token:      token,
		}, nil
	}
}

func (toktikServer *TokTikServer) CreateUserInfo(ctx context.Context, in *__.UserInfoRequest) (*__.UserBasicResponse, error) {

	return nil, nil
}
