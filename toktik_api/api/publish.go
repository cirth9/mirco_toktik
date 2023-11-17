package api

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"mime/multipart"
	"mirco_tiktok/toktik_api/forms"
	"mirco_tiktok/toktik_api/global"
	proto "mirco_tiktok/toktik_api/proto"
	"net/http"
	"strconv"
)

func fileToBytes(fileHeader *multipart.FileHeader) ([]byte, error) {
	size := fileHeader.Size
	bytesFile := make([]byte, size+1)
	file, err := fileHeader.Open()
	if err != nil {
		return []byte{}, err
	}
	_, err = file.Read(bytesFile)
	if err != nil {
		return []byte{}, err
	}
	return bytesFile, nil
}

func getVideoDataByte(fileHeader *multipart.FileHeader) ([]byte, error) {
	videoData := &forms.VideoData{}
	bytes, err := fileToBytes(fileHeader)
	if err != nil {
		return []byte{}, err
	}
	videoData.Content = bytes
	videoData.Filename = fileHeader.Filename
	videoData.Size = fileHeader.Size
	marshalData, err := json.Marshal(videoData)
	if err != nil {
		return []byte{}, err
	}
	return marshalData, nil
}

func PublishAction(c *gin.Context) {
	title := c.PostForm("title")
	token := c.PostForm("token")
	formData, err := c.MultipartForm()
	if err != nil {
		zap.S().Info(err)
		ResponseError(c, err)
		return
	}
	zap.S().Info("[PublishAction]", title, token)
	zap.S().Infof("%#v", formData)

	files := formData.File["data"]
	var publishActionResponse *proto.UserBasicResponse
	for _, v := range files {
		zap.S().Infof("%#v", v)
		videoData, err1 := getVideoDataByte(v)
		if err1 != nil {
			zap.S().Info(err1)
			ResponseError(c, err1)
			return
		}
		publishActionResponse, err1 = global.TokTikClient.PublishAction(context.Background(), &proto.UserPublishRequest{
			VideoData: videoData,
			Title:     title,
			Token:     token,
		})
		if err1 != nil {
			zap.S().Info(err1)
			ResponseError(c, err1)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": publishActionResponse.StatusCode,
		"status_msg":  publishActionResponse.StatusMsg,
	})
}

func UserPublishList(c *gin.Context) {
	//token := c.Query("token")
	userIdStr := c.Query("user_id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		zap.S().Info(err)
		ResponseError(c, err)
		return
	}
	UserPublishListResponse, err := global.TokTikClient.UserPublishList(context.Background(), &proto.UserPublishListRequest{
		UserId: userId,
	})
	if err != nil {
		zap.S().Info(err)
		ResponseError(c, err)
		return
	}

	var videoList []*forms.Video
	for _, v := range UserPublishListResponse.VideoInfoList {
		videoInfo := ProtoVideoInfoToFormsVideo(v)
		videoList = append(videoList, videoInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": UserPublishListResponse.UserBasicResponse.StatusCode,
		"status_msg":  UserPublishListResponse.UserBasicResponse.StatusMsg,
		"video_list":  videoList,
	})
}
