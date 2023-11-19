package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mirco_tiktok/toktik_srv/global"
	"mirco_tiktok/toktik_srv/models"
	proto "mirco_tiktok/toktik_srv/proto"
	"mirco_tiktok/toktik_srv/utils"
	"strconv"
	"strings"
	"time"
)

func (toktikServer *TokTikServer) UserRelationAction(ctx context.Context, in *proto.UserRelationRequest) (*proto.UserBasicResponse, error) {
	pool := goredis.NewPool(global.Rdb)
	rs := redsync.New(pool)

	parseToken, err := utils.ParseToken(in.Token)

	var userRelationFollowCache strings.Builder
	userRelationFollowCache.WriteString("relation-follow")
	userRelationFollowCache.WriteString(strconv.Itoa(int(parseToken.UserInfoID)))
	userRelationFollowCacheName := userRelationFollowCache.String()

	var userRelationFollowerCache strings.Builder
	userRelationFollowerCache.WriteString("relation-follower")
	userRelationFollowerCache.WriteString(strconv.Itoa(int(parseToken.UserInfoID)))
	userRelationFollowerCacheName := userRelationFollowerCache.String()

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
		RdbMutex := rs.NewMutex("relation-action")
		if err1 := RdbMutex.Lock(); err1 != nil {
			panic(err1)
		}
		_ = global.DB.Transaction(func(tx *gorm.DB) error {
			models.AddRelation(parseToken.UserInfoID, in.ToUserId)
			global.Rdb.Del(context.Background(), userRelationFollowCacheName)
			return nil
		})

		if ok, err1 := RdbMutex.Unlock(); !ok || err1 != nil {
			panic(err1)
		}

		return &proto.UserBasicResponse{
			StatusCode: 0,
			StatusMsg:  "关注成功",
			UserId:     int32(parseToken.UserInfoID),
			Token:      in.Token,
		}, nil
	} else if in.ActionType == "2" {
		RdbMutex := rs.NewMutex("relation-action")
		if err1 := RdbMutex.Lock(); err1 != nil {
			panic(err1)
		}

		_ = global.DB.Transaction(func(tx *gorm.DB) error {
			models.RemoveRelation(parseToken.UserInfoID, in.ToUserId)
			global.Rdb.Del(context.Background(), userRelationFollowerCacheName)
			return nil
		})

		if ok, err1 := RdbMutex.Unlock(); !ok || err1 != nil {
			panic(err1)
		}
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
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[UserRelationFollowList] Recover from panic!")
		}
	}()

	var userRelationCache strings.Builder
	userRelationCache.WriteString("relation-follow")
	userRelationCache.WriteString(strconv.Itoa(int(in.UserId)))
	userRelationCacheName := userRelationCache.String()

	if userRelationData, err := global.Rdb.Get(context.Background(), userRelationCacheName).Result(); err != nil {
		var userRelationCacheResponse *proto.UserRelationListResponse
		if err1 := json.Unmarshal([]byte(userRelationData), &userRelationCacheResponse); err1 != nil {
			panic(err1)
		}
		return userRelationCacheResponse, nil
	}

	UserFollowList := models.GetUserFollowList(in.UserId)
	zap.S().Info("follow list>>>>>>", UserFollowList)
	var UserInfoFollowList []*proto.UserInfo
	for _, v := range UserFollowList {
		userFollowInfo := UserToProtoUserInfo(v)
		UserInfoFollowList = append(UserInfoFollowList, userFollowInfo)
	}
	userRelationResponse := &proto.UserRelationListResponse{
		UserInfoList: UserInfoFollowList,
		UserBasicResponse: &proto.UserBasicResponse{
			StatusCode: 0,
			StatusMsg:  "get user follow list successfully",
			UserId:     int32(in.UserId),
			Token:      "",
		},
	}
	userRelationResponseData, err := json.Marshal(userRelationResponse)
	if err != nil {
		panic(err)
	}
	global.Rdb.Set(context.Background(), userRelationCacheName, userRelationResponseData, 0)
	return userRelationResponse, nil
}

func (toktikServer *TokTikServer) UserRelationFollowerList(ctx context.Context, in *proto.UserRelationFollowerListRequest) (*proto.UserRelationListResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[UserRelationFollowList] Recover from panic!")
		}
	}()

	var userRelationCache strings.Builder
	userRelationCache.WriteString("relation-follower")
	userRelationCache.WriteString(strconv.Itoa(int(in.UserId)))
	userRelationCacheName := userRelationCache.String()

	if userRelationData, err := global.Rdb.Get(context.Background(), userRelationCacheName).Result(); err != nil {
		var userRelationCacheResponse *proto.UserRelationListResponse
		if err1 := json.Unmarshal([]byte(userRelationData), &userRelationCacheResponse); err1 != nil {
			panic(err1)
		}
		return userRelationCacheResponse, nil
	}

	followerList := models.GetUserFollowerList(in.UserId)
	zap.S().Info("Follower list>>>>>>", followerList)
	var UserInfoFollowerList []*proto.UserInfo
	for _, v := range followerList {
		userFollowInfo := UserToProtoUserInfo(v)
		UserInfoFollowerList = append(UserInfoFollowerList, userFollowInfo)
	}
	userRelationResponse := &proto.UserRelationListResponse{
		UserInfoList: UserInfoFollowerList,
		UserBasicResponse: &proto.UserBasicResponse{
			StatusCode: 0,
			StatusMsg:  "get user follower list successfully",
			UserId:     int32(in.UserId),
			Token:      "",
		},
	}
	userRelationData, err := json.Marshal(userRelationResponse)
	if err != nil {
		panic(err)
	}
	global.Rdb.Set(context.Background(), userRelationCacheName, userRelationData, time.Hour)
	return userRelationResponse, nil
}
