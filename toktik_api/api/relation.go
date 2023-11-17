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

func UserRelationAction(c *gin.Context) {
	token := c.Query("token")
	toUserIdStr := c.Query("to_user_id")
	actionType := c.Query("action_type")
	toUserId, err := strconv.ParseInt(toUserIdStr, 10, 64)
	if err != nil {
		zap.S().Info(err)
		ResponseError(c, err)
		return
	}

	relationActionResponse, err := global.TokTikClient.UserRelationAction(context.Background(), &proto.UserRelationRequest{
		Token:      token,
		ToUserId:   toUserId,
		ActionType: actionType,
	})
	if err != nil {
		zap.S().Info(err)
		ResponseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": relationActionResponse.StatusCode,
		"status_msg":  relationActionResponse.StatusMsg,
	})
}

func UserRelationFollowList(c *gin.Context) {
	userIdStr := c.Query("user_id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		zap.S().Info(err)
		ResponseError(c, err)
		return
	}

	relationFollowListResponse, err := global.TokTikClient.UserRelationFollowList(context.Background(), &proto.UserRelationFollowListRequest{
		UserId: userId,
	})

	var userInfoList []*forms.UserInfo
	for _, v := range relationFollowListResponse.UserInfoList {
		userInfo := ProtoUserInfoToFormsUserInfo(v)
		userInfoList = append(userInfoList, userInfo)
	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": relationFollowListResponse.UserBasicResponse.StatusCode,
		"status_msg":  relationFollowListResponse.UserBasicResponse.StatusMsg,
		"user_list":   userInfoList,
	})
}

func UserRelationFollowerList(c *gin.Context) {
	userIdStr := c.Query("user_id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		zap.S().Info(err)
		ResponseError(c, err)
		return
	}

	relationFollowerListResponse, err := global.TokTikClient.UserRelationFollowerList(context.Background(), &proto.UserRelationFollowerListRequest{
		UserId: userId,
	})
	if err != nil {
		zap.S().Info(err)
		ResponseError(c, err)
		return
	}

	var userInfoList []*forms.UserInfo
	for _, v := range relationFollowerListResponse.UserInfoList {
		userInfo := ProtoUserInfoToFormsUserInfo(v)
		userInfoList = append(userInfoList, userInfo)
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": relationFollowerListResponse.UserBasicResponse.StatusCode,
		"status_msg":  relationFollowerListResponse.UserBasicResponse.StatusMsg,
		"user_list":   userInfoList,
	})
}
