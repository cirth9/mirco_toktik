package main

import (
	"flag"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"math"
	"mirco_tiktok/toktik_srv/config"
	"mirco_tiktok/toktik_srv/handler"
	"mirco_tiktok/toktik_srv/initialize"
	proto "mirco_tiktok/toktik_srv/proto"
	"mirco_tiktok/toktik_srv/utils"
	"mirco_tiktok/toktik_srv/utils/registery/consul"
	"net"
	"os"
	signal2 "os/signal"
	"syscall"
)

func main() {
	initialize.InitConfig()
	IP := flag.String("ip", config.TheServerConfig.Host, "ip地址")
	Port := flag.Int("port", 0, "port地址")
	flag.Parse()
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	toktikServer := grpc.NewServer(
		grpc.MaxSendMsgSize(math.MaxInt64),
		grpc.MaxRecvMsgSize(math.MaxInt64))
	proto.RegisterTokTikServer(toktikServer, &handler.TokTikServer{})

	//注册健康检查
	grpc_health_v1.RegisterHealthServer(toktikServer, health.NewServer())

	//服务注册
	consulRegistry := consul.NewRegistry(config.TheServerConfig.ConsulConfig.Host, config.TheServerConfig.ConsulConfig.Port)
	serviceId := uuid.NewV4().String()
	err2 := consulRegistry.Register(*IP, *Port, config.TheServerConfig.Tags, config.TheServerConfig.Name, serviceId)
	if err2 != nil {
		zap.S().Errorw("consul register error:", err2.Error())
	}

	listenAddress := fmt.Sprintf("%s:%d", *IP, *Port)
	zap.S().Info(listenAddress)

	lis, err := net.Listen("tcp", listenAddress)
	if err != nil {
		zap.S().Error(err)
		return
	}

	go func() {
		err = toktikServer.Serve(lis)
		if err != nil {
			zap.S().Error(err)
			return
		}
	}()

	signal := make(chan os.Signal)
	signal2.Notify(signal, syscall.SIGINT, syscall.SIGTERM)
	<-signal
	if err = consulRegistry.DeRegister(serviceId); err != nil {
		zap.S().Fatal("deregister failed")
	}
	zap.S().Info("deregister successfully")
}
