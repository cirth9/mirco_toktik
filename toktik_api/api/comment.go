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

func UserCommentAction(c *gin.Context) {
	token := c.Query("token")
	videoIdStr := c.Query("video_id")
	actionType := c.Query("action_type")
	commentText := c.Query("comment_text")
	commentId := c.Query("comment_id")
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		zap.S().Info(err)
		ResponseError(c, err)
		return
	}

	commentActionResponse, err := global.TokTikClient.UserCommentAction(context.Background(), &proto.UserCommentRequest{
		Token:       token,
		VideoId:     videoId,
		ActionType:  actionType,
		CommentText: commentText,
		CommentId:   commentId,
	})
	if err != nil {
		zap.S().Info(err)
		ResponseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": commentActionResponse.UserBasicResponse.StatusCode,
		"status_msg":  commentActionResponse.UserBasicResponse.StatusMsg,
		"comment":     ProtoCommentInfoToFormsCommentInfo(commentActionResponse.CommentInfo),
	})
}

func GetCommentList(c *gin.Context) {
	//token:=c.Query("token")
	videoIdStr := c.Query("video_id")
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		zap.S().Info(err)
		ResponseError(c, err)
		return
	}

	commentListResponse, err := global.TokTikClient.UserCommentList(context.Background(), &proto.UserCommentListRequest{
		VideoId: videoId,
	})
	if err != nil {
		zap.S().Info(err)
		ResponseError(c, err)
		return
	}

	var commentList []*forms.Comment
	for _, v := range commentListResponse.CommentInfoList {
		commentInfo := ProtoCommentInfoToFormsCommentInfo(v)
		commentList = append(commentList, commentInfo)
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code":  commentListResponse.UserBasicResponse.StatusCode,
		"status_msg":   commentListResponse.UserBasicResponse.StatusMsg,
		"comment_list": commentList,
	})
}
