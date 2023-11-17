package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mirco_tiktok/toktik_api/forms"
	"mirco_tiktok/toktik_api/global"
	proto "mirco_tiktok/toktik_api/proto"
	"net/http"
	"strconv"
)

func UserFavoriteAction(c *gin.Context) {
	token := c.Query("token")
	videoIdStr := c.Query("video_id")
	actionType := c.Query("action_type")

	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		zap.S().Info(err)
		ResponseError(c, err)
		return
	}

	favoriteActionResponse, err := global.TokTikClient.UserFavoriteAction(context.Background(), &proto.UserFavoriteRequest{
		Token:      token,
		VideoId:    videoId,
		ActionType: actionType,
	})
	c.JSON(http.StatusOK, gin.H{
		"status_code": favoriteActionResponse.StatusCode,
		"status_msg":  favoriteActionResponse.StatusMsg,
	})
}

func UserFavoriteList(c *gin.Context) {
	//token := c.Query("token")
	userIdStr := c.Query("user_id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		zap.S().Info(err)
		ResponseError(c, err)
		return
	}
	userFavoriteListResponse, err := global.TokTikClient.UserFavoriteList(context.Background(), &proto.UserFavoriteListRequest{
		UserId: userId,
	})
	if err != nil {
		zap.S().Info(err)
		ResponseError(c, err)
		return
	}

	var videoList []*forms.Video
	for _, v := range userFavoriteListResponse.VideoInfoList {
		videoInfo := ProtoVideoInfoToFormsVideo(v)
		videoList = append(videoList, videoInfo)
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": userFavoriteListResponse.UserBasicResponse.StatusCode,
		"status_msg":  userFavoriteListResponse.UserBasicResponse.StatusMsg,
		"video_list":  videoList,
	})
}
