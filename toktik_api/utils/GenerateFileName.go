package utils

import (
	"fmt"
	"time"
)

// GetVideoURL 生成video Url
func GetVideoURL(fileName string) string {
	//return fmt.Sprintf("http://%s:%s/%s/%s",
	//	config.TheServerConfig.ServerInfo.IP,
	//	config.TheServerConfig.ServerInfo.Port,
	//	config.TheServerConfig.StaticSavePath.DstName,
	//	fileName)
	return ""
}

// GetVideoFileName  生成独一无二的文件名
func GetVideoFileName(userId int64) string {
	return fmt.Sprintf("%d-%d", userId, time.Now().UnixNano())
}
