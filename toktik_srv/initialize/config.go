package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	goredislib "github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"mirco_tiktok/toktik_srv/config"
	"mirco_tiktok/toktik_srv/global"
	"os"
	"time"
)

var (
	groupId             string
	nacosConfigFileName string
)

func GetEnv(envName string, v *viper.Viper) bool {
	v.AutomaticEnv()
	return v.GetBool(envName)
}

func initMysql() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.TheServerConfig.MysqlConfig.UserName,
		config.TheServerConfig.MysqlConfig.Password,
		config.TheServerConfig.MysqlConfig.Host,
		config.TheServerConfig.MysqlConfig.Port,
		config.TheServerConfig.MysqlConfig.DBName)

	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             time.Second,
		Colorful:                  false,
		IgnoreRecordNotFoundError: false,
		ParameterizedQueries:      false,
		LogLevel:                  logger.Info,
	})

	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Println(">>>>>>", err)
		return
	}
	return err
}

func initRedis() {
	addr := fmt.Sprintf("%s:%d",
		config.TheServerConfig.RedisConfig.IP,
		config.TheServerConfig.RedisConfig.Port)
	global.Rdb = goredislib.NewClient(&goredislib.Options{
		Addr:     addr,
		Password: "qq31415926535--",
		DB:       0,
	})
}

func InitConfig() {
	err := InitLogger()
	if err != nil {
		log.Println(err)
	}

	v := viper.New()
	if GetEnv("TIKTOK_ENV", v) {
		groupId = "dev"
	} else {
		groupId = "pro"
	}
	zap.S().Info(groupId)

	nacosConfigFileName = "nacos-config.yml"
	v.SetConfigFile(nacosConfigFileName)
	if err1 := v.ReadInConfig(); err1 != nil {
		zap.S().Panic(err1)
	}

	err2 := v.Unmarshal(&config.TheNacosConfig)
	if err2 != nil {
		zap.S().Panic(err2)
	}

	log.Printf("%+v", config.TheNacosConfig)
	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: config.TheNacosConfig.NacosServer.Ip,
			Port:   config.TheNacosConfig.NacosServer.Port,
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         config.TheNacosConfig.NacosClient.NamespaceId, // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           config.TheNacosConfig.NacosClient.TimeoutMs,
		NotLoadCacheAtStart: config.TheNacosConfig.NacosClient.NotLoadCacheAtStart,
		LogDir:              config.TheNacosConfig.NacosClient.LogDir,
		CacheDir:            config.TheNacosConfig.NacosClient.CacheDir,
	}

	// 创建动态配置客户端
	client, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(err)
	}
	userSrvConfig, err := client.GetConfig(vo.ConfigParam{
		DataId: config.TheNacosConfig.NacosServer.DataId,
		Group:  groupId,
	})
	log.Println(">>>>>>", userSrvConfig)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(userSrvConfig), &config.TheServerConfig)
	if err != nil {
		zap.S().Panic(err)
	}

	zap.S().Info("TheServerConfig: ", config.TheServerConfig)

	err = initMysql()
	if err != nil {
		zap.S().Error(err)
	}
	initRedis()
}
