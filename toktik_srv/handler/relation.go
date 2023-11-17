package handler

import (
	"context"
	"go.uber.org/zap"
	"mirco_tiktok/toktik_srv/models"
	proto "mirco_tiktok/toktik_srv/proto"
	"mirco_tiktok/toktik_srv/utils"
)

func (toktikServer *TokTikServer) UserRelationAction(ctx context.Context, in *proto.UserRelationRequest) (*proto.UserBasicResponse, error) {
	parseToken, err := utils.ParseToken(in.Token)
	if err != nil {
		zap.S().Error(err)
		return &proto.UserBasicResponse{
			StatusCode: -1,
			StatusMsg:  "token parse error",
			UserId:     0,
			Token:      in.Token,
		}, err
	}
	if in.ActionType == "1" {
		models.AddRelation(parseToken.UserInfoID, in.ToUserId)
		return &proto.UserBasicResponse{
			StatusCode: 0,
			StatusMsg:  "关注成功",
			UserId:     int32(parseToken.UserInfoID),
			Token:      in.Token,
		}, nil
	} else if in.ActionType == "2" {
		models.RemoveRelation(parseToken.UserInfoID, in.ToUserId)
		return &proto.UserBasicResponse{
			StatusCode: 0,
			StatusMsg:  "取消关注成功",
			UserId:     int32(parseToken.UserInfoID),
			Token:      in.Token,
		}, nil
	} else {
		return &proto.UserBasicResponse{
			StatusCode: 0,
			StatusMsg:  "未知操作",
			UserId:     int32(parseToken.UserInfoID),
			Token:      in.Token,
		}, err
	}
}

func (toktikServer *TokTikServer) UserRelationFollowList(ctx context.Context, in *proto.UserRelationFollowListRequest) (*proto.UserRelationListResponse, error) {
	UserFollowList := models.GetUserFollowList(in.UserId)
	zap.S().Info("follow list>>>>>>", UserFollowList)
	var UserInfoFollowList []*proto.UserInfo
	for _, v := range UserFollowList {
		userFollowInfo := UserToProtoUserInfo(v)
		UserInfoFollowList = append(UserInfoFollowList, userFollowInfo)
	}
	return &proto.UserRelationListResponse{
		UserInfoList: UserInfoFollowList,
		UserBasicResponse: &proto.UserBasicResponse{
			StatusCode: 0,
			StatusMsg:  "get user follow list successfully",
			UserId:     int32(in.UserId),
			Token:      "",
		},
	}, nil
}

func (toktikServer *TokTikServer) UserRelationFollowerList(ctx context.Context, in *proto.UserRelationFollowerListRequest) (*proto.UserRelationListResponse, error) {
	followerList := models.GetUserFollowerList(in.UserId)
	zap.S().Info("Follower list>>>>>>", followerList)
	var UserInfoFollowerList []*proto.UserInfo
	for _, v := range followerList {
		userFollowInfo := UserToProtoUserInfo(v)
		UserInfoFollowerList = append(UserInfoFollowerList, userFollowInfo)
	}
	return &proto.UserRelationListResponse{
		UserInfoList: UserInfoFollowerList,
		UserBasicResponse: &proto.UserBasicResponse{
			StatusCode: 0,
			StatusMsg:  "get user follower list successfully",
			UserId:     int32(in.UserId),
			Token:      "",
		},
	}, nil
}
