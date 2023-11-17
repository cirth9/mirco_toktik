package router

import (
	"github.com/gin-gonic/gin"

	"mirco_tiktok/toktik_api/api"

	"mirco_tiktok/toktik_api/middleware"
)

func InitTokTikRouter(r *gin.RouterGroup) {

	r.GET("/feed", api.VideoFeed)

	userGroup := r.Group("user")
	{
		//middleware.AdminAuth()
		userGroup.GET("/", middleware.JwtMiddleware(), api.GetUserInfo)
		userGroup.POST("/login/", api.UserLogin)
		userGroup.POST("/register/", api.UserRegister)
	}

	publishGroup := r.Group("publish")
	{
		publishGroup.POST("/action/", middleware.JwtMiddleware(), api.PublishAction)
		publishGroup.GET("/list/", middleware.JwtMiddleware(), api.UserPublishList)
	}

	favorGroup := r.Group("favorite")
	{
		favorGroup.POST("/action/", middleware.JwtMiddleware(), api.UserFavoriteAction)
		favorGroup.GET("/list/", middleware.JwtMiddleware(), api.UserFavoriteList)
	}

	commentGroup := r.Group("comment")
	{
		commentGroup.POST("/action/", middleware.JwtMiddleware(), api.UserCommentAction)
		commentGroup.GET("/list/", middleware.JwtMiddleware(), api.GetCommentList)
	}

	relationGroup := r.Group("relation")
	{
		relationGroup.POST("/action/", middleware.JwtMiddleware(), api.UserRelationAction)
		relationGroup.GET("/follow/list/", middleware.JwtMiddleware(), api.UserRelationFollowList)
		relationGroup.GET("/follower/list/", middleware.JwtMiddleware(), api.UserRelationFollowerList)
	}

	//后面的东西前端都没有实现，那我也没办法了
	r.GET("/relation/friend/list/", nil)

	r.POST("/message/chat/", nil)
	r.GET("/message/action/", nil)
}
