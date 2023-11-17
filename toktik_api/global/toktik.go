package global

import (
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
	proto "mirco_tiktok/toktik_api/proto"
	"time"
)

// TokenExpireDuration token过期时间
const TokenExpireDuration = time.Hour * 2

// Secret 密钥
var Secret = []byte("secret")

// Rdb redis client
var (
	Rdb         *redis.Client
	ExpiredTime = time.Minute * 2
)

var (
	//GrpcAddress GRPC地址
	GrpcAddress string

	TokTikConn   *grpc.ClientConn
	TokTikClient proto.TokTikClient
)
