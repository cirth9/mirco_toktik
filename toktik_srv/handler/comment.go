package handler

import (
	"context"
	"encoding/json"
	"fmt"
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
	"time"
)

func (toktikServer *TokTikServer) UserCommentAction(ctx context.Context, in *proto.UserCommentRequest) (*proto.UserCommentResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[UserCommentAction] Recover from panic!")
		}
	}()

	pool := goredis.NewPool(global.Rdb)
	rs := redsync.New(pool)
	rdsMutex := rs.NewMutex("comment_count")

	var userCommentCache strings.Builder
	userCommentCache.WriteString("comment-")
	userCommentCache.WriteString(strconv.Itoa(int(in.VideoId)))
	userCommentCacheName := userCommentCache.String()

	parseToken, err := utils.ParseToken(in.Token)
	if err != nil {
		zap.S().Info(err)
		return &proto.UserCommentResponse{
			CommentInfo: nil,
			UserBasicResponse: &proto.UserBasicResponse{
				StatusCode: -1,
				StatusMsg:  "评论操作失败",
				UserId:     0,
				Token:      in.Token,
			},
		}, err
	}

	newComment := models2.Comment{
		UserInfoId: parseToken.UserInfoID,
		VideoId:    in.VideoId,
		User:       models2.FindUserInfoById(parseToken.UserInfoID),
		Content:    in.CommentText,
		CreatedAt:  time.Time{},
		CreateDate: time.Now().String(),
	}

	if in.ActionType == "1" {

		if err = rdsMutex.Lock(); err != nil {
			panic(err)
		}

		_ = global.DB.Transaction(func(tx *gorm.DB) error {
			models2.AddNewComment(newComment)
			models2.AscCommentCountById(in.VideoId)
			global.Rdb.Del(context.Background(), userCommentCacheName)
			return nil
		})

		if ok, err1 := rdsMutex.Unlock(); !ok || err1 != nil {
			panic(err1)
		}

		return &proto.UserCommentResponse{
			CommentInfo: CommentToProtoCommentInfo(&newComment),
			UserBasicResponse: &proto.UserBasicResponse{
				StatusCode: 0,
				StatusMsg:  "评论成功",
				UserId:     int32(parseToken.UserInfoID),
				Token:      in.Token,
			},
		}, nil
	} else if in.ActionType == "2" {
		CommentIdInt, err1 := strconv.ParseInt(in.CommentId, 10, 64)
		if err1 != nil {
			zap.S().Info(err1)
			return &proto.UserCommentResponse{
				CommentInfo: nil,
				UserBasicResponse: &proto.UserBasicResponse{
					StatusCode: -1,
					StatusMsg:  "删除评论失败",
					UserId:     0,
					Token:      in.Token,
				},
			}, err1
		}

		if err = rdsMutex.Lock(); err != nil {
			panic(err)
		}

		_ = global.DB.Transaction(func(tx *gorm.DB) error {
			models2.RemoveComment(CommentIdInt)
			models2.DscCommentCountById(in.VideoId)
			global.Rdb.Del(context.Background(), userCommentCacheName)
			return nil
		})

		if ok, err2 := rdsMutex.Unlock(); !ok || err2 != nil {
			panic(err2)
		}

		return &proto.UserCommentResponse{
			CommentInfo: nil,
			UserBasicResponse: &proto.UserBasicResponse{
				StatusCode: 0,
				StatusMsg:  "删除评论成功",
				UserId:     int32(parseToken.UserInfoID),
				Token:      in.Token,
			},
		}, nil
	} else {
		return &proto.UserCommentResponse{
			CommentInfo: nil,
			UserBasicResponse: &proto.UserBasicResponse{
				StatusCode: -1,
				StatusMsg:  "操作类型不存在",
				UserId:     0,
				Token:      in.Token,
			},
		}, err
	}
}

func (toktikServer *TokTikServer) UserCommentList(ctx context.Context, in *proto.UserCommentListRequest) (*proto.UserCommentListResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[UserFavoriteList] Recover from panic!")
		}
	}()

	var userCommentCache strings.Builder
	userCommentCache.WriteString("comment-")
	userCommentCache.WriteString(strconv.Itoa(int(in.VideoId)))
	userCommentCacheName := userCommentCache.String()
	if userCommentCacheData, err := global.Rdb.Get(context.Background(), userCommentCacheName).Result(); err != nil {
		var userCommentCacheResponse *proto.UserCommentListResponse
		if err1 := json.Unmarshal([]byte(userCommentCacheData), &userCommentCacheResponse); err1 != nil {
			panic(err1)
		}
		return userCommentCacheResponse, nil
	}

	commentList := models2.GetAllCommentByVid(in.VideoId)
	var commentInfoList []*proto.CommentInfo
	for _, v := range commentList {
		commentInfo := CommentToProtoCommentInfo(v)
		commentInfoList = append(commentInfoList, commentInfo)
	}
	userCommentListResponse := &proto.UserCommentListResponse{
		CommentInfoList: commentInfoList,
		UserBasicResponse: &proto.UserBasicResponse{
			StatusCode: 0,
			StatusMsg:  "获取评论列表成功",
		},
	}
	userCommentData, err := json.Marshal(userCommentListResponse)
	if err != nil {
		panic(err)
	}
	global.Rdb.Set(context.Background(), userCommentCacheName, userCommentData, time.Hour)
	return userCommentListResponse, nil
}
