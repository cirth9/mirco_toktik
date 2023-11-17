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

func VideoFeed(c *gin.Context) {
	token := c.Query("token")
	latestTimeStr := c.Query("latest_time")
	latestTime, err := strconv.ParseInt(latestTimeStr, 10, 64)
	if err != nil {
		zap.S().Info(err)
		ResponseError(c, err)
		return
	}
	VideoFeedResponse, err := global.TokTikClient.VideoFeed(context.Background(), &proto.VideoFeedRequest{
		Token:      token,
		LatestTime: latestTime,
	})
	if err != nil {
		zap.S().Info(err)
		ResponseError(c, err)
		return
	}

	var videoList []*forms.Video
	for _, v := range VideoFeedResponse.VideoList {
		videoInfo := ProtoVideoInfoToFormsVideo(v)
		videoList = append(videoList, videoInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": VideoFeedResponse.UserBasicResponse.StatusCode,
		"status_msg":  VideoFeedResponse.UserBasicResponse.StatusMsg,
		"next_time":   VideoFeedResponse.NextTime,
		"video_list":  videoList,
	})
}
