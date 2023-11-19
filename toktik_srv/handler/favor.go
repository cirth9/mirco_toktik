package handler

import (
	"context"
	"encoding/json"
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
	"strconv"
	"strings"
)

// UserFavoriteAction 这里需要考虑并发问题
func (toktikServer *TokTikServer) UserFavoriteAction(ctx context.Context, in *proto.UserFavoriteRequest) (*proto.UserBasicResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[UserFavoriteAction] Recover from panic!")
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

	var userFavorCache strings.Builder
	userFavorCache.WriteString("favor-")
	userFavorCache.WriteString(strconv.Itoa(int(parseToken.UserInfoID)))
	userFavorCacheName := userFavorCache.String()

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
			//global.Rdb.SAdd(context.Background(), userFavorCache.String(), strconv.Itoa(int(in.VideoId)))
			global.Rdb.Del(context.Background(), userFavorCacheName)
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
			//global.Rdb.SRem(context.Background(), userFavorCache.String(), strconv.Itoa(int(in.VideoId)))
			global.Rdb.Del(context.Background(), userFavorCacheName)
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
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[UserFavoriteList] Recover from panic!")
		}
	}()

	var userFavorCache strings.Builder
	userFavorCache.WriteString("favor-")
	userFavorCache.WriteString(strconv.Itoa(int(in.UserId)))
	userFavorCacheKey := userFavorCache.String()
	if favorListCache, err := global.Rdb.Get(context.Background(), userFavorCacheKey).Result(); err != nil {
		var favorListData *proto.UserFavoriteListResponse
		err1 := json.Unmarshal([]byte(favorListCache), &favorListData)
		if err1 != nil {
			panic(err1)
		}
		return favorListData, nil
	}

	videoList := models2.GetUserFavorList(in.UserId)
	var videoInfoList []*proto.VideoInfo
	for _, v := range videoList {
		videoInfo := VideoToProtoVideoInfo(v)
		videoInfoList = append(videoInfoList, videoInfo)
	}
	userFavorListResponse := &proto.UserFavoriteListResponse{
		VideoInfoList: videoInfoList,
		UserBasicResponse: &proto.UserBasicResponse{
			StatusCode: 0,
			StatusMsg:  "get user favor list successfully",
			UserId:     int32(in.UserId),
			Token:      "",
		},
	}
	userFavorListData, err := json.Marshal(userFavorListResponse)
	if err != nil {
		panic(err)
	}

	global.Rdb.Set(context.Background(), userFavorCacheKey, userFavorListData, 0)
	return userFavorListResponse, nil
}
