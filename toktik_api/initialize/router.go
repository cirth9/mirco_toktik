package initialize

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mirco_tiktok/toktik_api/config"
	"mirco_tiktok/toktik_api/middleware"
	"mirco_tiktok/toktik_api/router"
)

func InitRouters() *gin.Engine {
	r := gin.Default()
	r.Static(config.TheServerConfig.StaticSavePath.DstName, config.TheServerConfig.StaticSavePath.DST)
	zap.S().Info(config.TheServerConfig.StaticSavePath.DST)
	r.Use(middleware.Cors())
	userGroup := r.Group("/douyin")
	router.InitTokTikRouter(userGroup)
	return r
}
