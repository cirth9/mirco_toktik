package handler

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
	"mirco_tiktok/toktik_srv/models"
	proto "mirco_tiktok/toktik_srv/proto"
	"mirco_tiktok/toktik_srv/utils"
	"time"
)

func (toktikServer *TokTikServer) VideoFeed(ctx context.Context, in *proto.VideoFeedRequest) (*proto.VideoFeedResponse, error) {
	if in.Token == "" {
		return DoWithoutToken(in.LatestTime)
	} else {
		return DoWithToken(in.LatestTime, in.Token)
	}
}

func DoWithToken(latestTime int64, token string) (*proto.VideoFeedResponse, error) {
	//客制化内容，随便写点东西吧
	//这里就随便写个，无法刷到自己投递的内容
	TimeStampInt := latestTime
	zap.S().Info("TimeStampInt >>>>>", TimeStampInt)

	parseToken, err := utils.ParseToken(token)
	if err != nil {
		zap.S().Error(err)
		return &proto.VideoFeedResponse{
			VideoList: nil,
			UserBasicResponse: &proto.UserBasicResponse{
				StatusCode: -1,
				StatusMsg:  err.Error(),
				UserId:     int32(parseToken.UserInfoID),
				Token:      "",
			},
			NextTime: 0,
		}, nil
	}
	log.Printf("%+v", parseToken)
	videoList := models.GetVideoListWithoutSelf(time.Unix(0, TimeStampInt*1e6), parseToken.UserInfoID)
	fmt.Printf(" videoList >>>>> %+v\n", videoList)

	var nextTime int64
	if len(videoList) == 0 {
		nextTime = 0
	} else {
		nextTime = videoList[0].CreatedAt.Unix()
	}

	var videoInfoList []*proto.VideoInfo
	for _, v := range videoList {
		videoInfo := VideoToProtoVideoInfo(v)
		videoInfoList = append(videoInfoList, videoInfo)
	}

	return &proto.VideoFeedResponse{
		VideoList: videoInfoList,
		UserBasicResponse: &proto.UserBasicResponse{
			StatusCode: 0,
			StatusMsg:  "get user's publish-list successfully",
			UserId:     int32(parseToken.UserInfoID),
			Token:      "",
		},
		NextTime: nextTime,
	}, nil
}

func DoWithoutToken(latestTime int64) (*proto.VideoFeedResponse, error) {
	TimeStampInt := latestTime
	zap.S().Info("TimeStampInt >>>>>", TimeStampInt)

	videoList := models.GetVideoList(time.Unix(0, TimeStampInt*1e6))
	fmt.Printf(" videoList >>>>> %+v\n", videoList)

	var nextTime int64
	if len(videoList) == 0 {
		nextTime = 0
	} else {
		nextTime = videoList[0].CreatedAt.Unix()
	}

	var videoInfoList []*proto.VideoInfo
	for _, v := range videoList {
		videoInfo := VideoToProtoVideoInfo(v)
		videoInfoList = append(videoInfoList, videoInfo)
	}

	return &proto.VideoFeedResponse{
		VideoList: videoInfoList,
		UserBasicResponse: &proto.UserBasicResponse{
			StatusCode: 0,
			StatusMsg:  "get user's publish-list successfully",
			UserId:     0,
			Token:      "",
		},
		NextTime: nextTime,
	}, nil
}
