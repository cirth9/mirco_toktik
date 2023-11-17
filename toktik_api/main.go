package main

import (
	"fmt"
	"go.uber.org/zap"
	"mirco_tiktok/toktik_api/initialize"
	"mirco_tiktok/toktik_api/utils"
)

func main() {
	var err error
	//初始化全局zap
	initialize.InitLogger()

	initialize.InitConfig()

	r := initialize.InitRouters()

	port, err := utils.GetFreePort()

	port = 9090

	zap.S().Info("server start,port :", port)
	addr := fmt.Sprintf(":%d", port)
	err = r.Run(addr)
	if err != nil {
		zap.S().Panic("start failed", err.Error())
	}
}
