package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mirco_tiktok/toktik_api/forms"
	"mirco_tiktok/toktik_api/global"
	proto "mirco_tiktok/toktik_api/proto"
	"net/http"
	"strconv"
)

func GrpcErrorToHTTP(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误",
				})
			}
		}
	}
}

func UserInfoResponseToUserInfo(userInfoResponse *proto.UserInfoResponse) *forms.UserInfo {
	return &forms.UserInfo{
		Id:              userInfoResponse.Id,
		Name:            userInfoResponse.Name,
		FollowCount:     userInfoResponse.FollowCount,
		FollowerCount:   userInfoResponse.FollowerCount,
		IsFollow:        userInfoResponse.IsFollow,
		Avatar:          userInfoResponse.Avatar,
		BackgroundImage: userInfoResponse.BackgroundImage,
		Signature:       userInfoResponse.Signature,
		TotalFavorited:  userInfoResponse.TotalFavorited,
		WorkCount:       userInfoResponse.WorkCount,
		FavoriteCount:   userInfoResponse.FavoriteCount,
	}
}

func GetUserInfo(c *gin.Context) {
	userIdStr := c.Query("user_id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		zap.S().Info(err)
		ResponseError(c, err)
		return
	}

	UserInfoResponse, err := global.TokTikClient.GetUserInfo(context.Background(), &proto.UserInfoRequest{
		UserId: userId,
	})
	if err != nil {
		zap.S().Info(err)
		ResponseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": UserInfoResponse.UserBasicResponse.StatusCode,
		"status_msg":  UserInfoResponse.UserBasicResponse.StatusMsg,
		"user":        UserInfoResponseToUserInfo(UserInfoResponse),
	})
}

func UserLogin(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	UserLoginResponse, err := global.TokTikClient.UserLogin(context.Background(), &proto.UserInfoRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		zap.S().Info(err)
		ResponseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": UserLoginResponse.StatusCode,
		"status_msg":  UserLoginResponse.StatusMsg,
		"user_id":     UserLoginResponse.UserId,
		"token":       UserLoginResponse.Token,
	})
}

func UserRegister(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	UserRegisterResponse, err := global.TokTikClient.UserRegister(context.Background(), &proto.UserInfoRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		zap.S().Info(err)
		ResponseError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": UserRegisterResponse.StatusCode,
		"status_msg":  UserRegisterResponse.StatusMsg,
		"user_id":     UserRegisterResponse.UserId,
		"token":       UserRegisterResponse.Token,
	})
}

func ResponseError(c *gin.Context, err error) {
	c.JSON(http.StatusAccepted, gin.H{
		"status_code": -1,
		"status_msg":  "StatusInternalServerError",
		"error":       err.Error(),
	})
}
