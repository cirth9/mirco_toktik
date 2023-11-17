package handler

import (
	"context"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mirco_tiktok/toktik_srv/global"
	models2 "mirco_tiktok/toktik_srv/models"
	proto "mirco_tiktok/toktik_srv/proto"
	"mirco_tiktok/toktik_srv/utils"
)

// UserFavoriteAction 这里需要考虑并发问题
func (toktikServer *TokTikServer) UserFavoriteAction(ctx context.Context, in *proto.UserFavoriteRequest) (*proto.UserBasicResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[UserFavorite] Recover from panic!")
		}
	}()

	pool := goredis.NewPool(global.Rdb)
	rs := redsync.New(pool)

	zap.S().Info(in)
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
		if models2.IsUserFavorVideo(parseToken.UserInfoID, in.VideoId) {
			return &proto.UserBasicResponse{
				StatusCode: -1,
				StatusMsg:  "favor repeat",
				UserId:     int32(parseToken.UserInfoID),
				Token:      in.Token,
			}, errors.Errorf("favor repeat")
		}
		rdsMutex := rs.NewMutex("favor_count")
		if err = rdsMutex.Lock(); err != nil {
			panic(err)
		}

		_ = global.DB.Transaction(func(tx *gorm.DB) error {
			models2.AscFavoriteCountById(in.VideoId)
			models2.AddUserFavorVideos(parseToken.UserInfoID, in.VideoId)
			return nil
		})

		if ok, err1 := rdsMutex.Unlock(); !ok || err1 != nil {
			panic(err1)
		}

		return &proto.UserBasicResponse{
			StatusCode: 0,
			StatusMsg:  "favor successfully",
			UserId:     int32(parseToken.UserInfoID),
			Token:      in.Token,
		}, nil
	} else if in.ActionType == "2" {
		if !models2.IsUserFavorVideo(parseToken.UserInfoID, in.VideoId) {
			return &proto.UserBasicResponse{
				StatusCode: -1,
				StatusMsg:  "favor quit repeat",
				UserId:     int32(parseToken.UserInfoID),
				Token:      in.Token,
			}, errors.Errorf("favor quit repeat")
		}

		rdsMutex := rs.NewMutex("favor_count")
		if err = rdsMutex.Lock(); err != nil {
			panic(err)
		}

		_ = global.DB.Transaction(func(tx *gorm.DB) error {
			models2.DscFavoriteCountById(in.VideoId)
			models2.RmUserFavorVideos(parseToken.UserInfoID, in.VideoId)
			return nil
		})

		if ok, err1 := rdsMutex.Unlock(); !ok || err1 != nil {
			panic(err1)
		}

		return &proto.UserBasicResponse{
			StatusCode: 0,
			StatusMsg:  "favor quit",
			UserId:     int32(parseToken.UserInfoID),
			Token:      in.Token,
		}, nil
	} else {
		return &proto.UserBasicResponse{
			StatusCode: -1,
			StatusMsg:  "action_type undefined",
			UserId:     int32(parseToken.UserInfoID),
			Token:      in.Token,
		}, errors.Errorf("action_type undefined")
	}
}

func (toktikServer *TokTikServer) UserFavoriteList(ctx context.Context, in *proto.UserFavoriteListRequest) (*proto.UserFavoriteListResponse, error) {
	videoList := models2.GetUserFavorList(in.UserId)
	var videoInfoList []*proto.VideoInfo
	for _, v := range videoList {
		videoInfo := VideoToProtoVideoInfo(v)
		videoInfoList = append(videoInfoList, videoInfo)
	}
	return &proto.UserFavoriteListResponse{
		VideoInfoList: videoInfoList,
		UserBasicResponse: &proto.UserBasicResponse{
			StatusCode: 0,
			StatusMsg:  "get user favor list successfully",
			UserId:     int32(in.UserId),
			Token:      "",
		},
	}, nil
}
