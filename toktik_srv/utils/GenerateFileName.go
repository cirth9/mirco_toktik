package utils

import (
	"fmt"
	"mirco_tiktok/toktik_srv/config"
	"time"
)

// GetVideoURL 生成video Url
func GetVideoURL(fileName string) string {
	return fmt.Sprintf("http://%s:%d/%s/%s",
		config.TheServerConfig.ServerInfo.IP,
		config.TheServerConfig.ServerInfo.Port,
		config.TheServerConfig.StaticSavePath.DstName,
		fileName)
}

// GetVideoFileName  生成独一无二的文件名
func GetVideoFileName(userId int64) string {
	return fmt.Sprintf("%d-%d", userId, time.Now().UnixNano())
}
