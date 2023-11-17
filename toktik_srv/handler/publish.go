package handler

import (
	"context"
	"encoding/json"
	"github.com/go-errors/errors"
	"go.uber.org/zap"
	"mirco_tiktok/toktik_srv/config"
	models2 "mirco_tiktok/toktik_srv/models"
	proto "mirco_tiktok/toktik_srv/proto"
	utils2 "mirco_tiktok/toktik_srv/utils"
	"os"
	"path/filepath"
)

var (
	VideoFormat = map[string]interface{}{
		".mp4": nil,
		".avi": nil,
		".mov": nil,
		".wmv": nil,
		".flv": nil,
	}
)

type VideoData struct {
	Content  []byte `json:"content"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
}

func saveFile(file *VideoData, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = out.Write(file.Content)
	if err != nil {
		return err
	}
	return nil
}

func (toktikServer *TokTikServer) PublishAction(ctx context.Context, in *proto.UserPublishRequest) (*proto.UserBasicResponse, error) {
	var parseToken *utils2.UserBasicClaims
	parseToken, err := utils2.ParseToken(in.Token)
	if err != nil {
		zap.S().Error(err)
		return &proto.UserBasicResponse{
			StatusCode: -1,
			StatusMsg:  "publish failed",
			UserId:     0,
			Token:      in.Token,
		}, errors.Errorf("publish failed,error: ", err.Error())
	}

	var videoData *VideoData
	err = json.Unmarshal(in.VideoData, &videoData)
	if err != nil {
		zap.S().Error("[PublishAction] json unmarshal error:", err)
		return &proto.UserBasicResponse{
			StatusCode: -1,
			StatusMsg:  "unmarshal error",
			UserId:     0,
			Token:      "",
		}, errors.Errorf("json unmarshal error:", err.Error())
	}

	suffix := filepath.Ext(videoData.Filename)
	if _, ok := VideoFormat[suffix]; !ok {
		zap.S().Info("不支持的视频格式")
		return &proto.UserBasicResponse{
			StatusCode: -1,
			StatusMsg:  "不支持的视频格式",
			UserId:     int32(parseToken.UserInfoID),
			Token:      in.Token,
		}, errors.Errorf("不支持的视频格式")
	}

	filename := utils2.GetVideoFileName(parseToken.UserInfoID) + suffix
	savePath := filepath.Join(config.TheServerConfig.StaticSavePath.Dst, filename)
	zap.S().Info(savePath)

	err = saveFile(videoData, savePath)
	if err != nil {
		zap.S().Error(err)
		return &proto.UserBasicResponse{
			StatusCode: -1,
			StatusMsg:  "publish failed",
			UserId:     0,
			Token:      in.Token,
		}, errors.Errorf("publish failed")
	}

	PublishVideo := models2.Video{
		UserInfoId: parseToken.UserInfoID,
		Author:     models2.FindUserInfoById(parseToken.UserInfoID),
		PlayUrl:    utils2.GetVideoURL(filename),
		CoverUrl:   "",
		Title:      in.Title,
	}
	models2.CreateVideo(PublishVideo)
	zap.S().Infof("%+v", PublishVideo.Author)
	if err != nil {
		zap.S().Error("save file error", err)
		return &proto.UserBasicResponse{
			StatusCode: -1,
			StatusMsg:  "publish failed",
			UserId:     int32(parseToken.UserInfoID),
			Token:      in.Token,
		}, errors.Errorf("publish failed")
	}

	return &proto.UserBasicResponse{
		StatusCode: 0,
		StatusMsg:  "publish successfully",
		UserId:     int32(parseToken.UserInfoID),
		Token:      in.Token,
	}, nil
}

func (toktikServer *TokTikServer) UserPublishList(ctx context.Context, in *proto.UserPublishListRequest) (*proto.UserPublishListResponse, error) {
	videoList := models2.FindVideoByUid(in.UserId)
	var videoInfoList []*proto.VideoInfo
	for _, v := range videoList {
		videoInfo := VideoToProtoVideoInfo(v)
		videoInfoList = append(videoInfoList, videoInfo)
	}
	return &proto.UserPublishListResponse{
		VideoInfoList: videoInfoList,
		UserBasicResponse: &proto.UserBasicResponse{
			StatusCode: 0,
			StatusMsg:  "get user publish list successfully",
			UserId:     int32(in.UserId),
			Token:      "",
		},
	}, nil
}
