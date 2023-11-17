package config

type MysqlConfig struct {
	UserName string `yaml:"username" json:"username"  mapstructure:"username"`
	Password string `yaml:"password" json:"password"  mapstructure:"password"`
	Host     string `yaml:"host" json:"host"  mapstructure:"host"`
	Port     int    `yaml:"port" json:"port"  mapstructure:"port"`
	DBName   string `yaml:"dbname" json:"dbname"  mapstructure:"dbname"`
}

type ConsulConfig struct {
	Host string `yaml:"host" json:"host"  mapstructure:"host"`
	Port int    `yaml:"port" json:"port"  mapstructure:"port"`
}

type ServerInfo struct {
	IP   string `yaml:"ip" json:"ip"  mapstructure:"ip"`
	Port int    `yaml:"port" json:"port"  mapstructure:"port"`
}

type RedisConfig struct {
	Password string `yaml:"password" json:"password"  mapstructure:"password"`
	IP       string `yaml:"ip" json:"ip"  mapstructure:"ip"`
	Port     int    `yaml:"port" json:"port"  mapstructure:"port"`
	Db       int    `yaml:"db" json:"db"  mapstructure:"db"`
}

type StaticSavePath struct {
	Dst     string `yaml:"dst" json:"dst"  mapstructure:"dst"`
	DstName string `yaml:"dst_name" json:"dst_name"  mapstructure:"dst_name"`
}

type NacosConfig struct {
	NacosServer NacosServer `yaml:"nacos_server" mapstructure:"nacos_server"`
	NacosClient NacosClient `yaml:"nacos_client" mapstructure:"nacos_client"`
}

type NacosServer struct {
	DataId string `yaml:"dataId" mapstructure:"dataId"`
	Ip     string `yaml:"ip"mapstructure:"ip"`
	Port   uint64 `yaml:"port"mapstructure:"port"`
}

type NacosClient struct {
	NotLoadCacheAtStart bool   `yaml:"not_load_cache_at_start" mapstructure:"not_load_cache_at_start"`
	LogDir              string `yaml:"log_dir" mapstructure:"log_dir"`
	CacheDir            string `yaml:"cache_dir" mapstructure:"cache_dir"`
	NamespaceId         string `yaml:"namespace_id" mapstructure:"namespace_id"`
	TimeoutMs           uint64 `yaml:"timeout_ms" mapstructure:"timeout_ms"`
}

type ServerConfig struct {
	Name           string   `yaml:"name" json:"name"  mapstructure:"name"`
	Host           string   `yaml:"host" json:"host"  mapstructure:"host"`
	Tags           []string `yaml:"tags" json:"tags"  mapstructure:"tags"`
	MysqlConfig    `yaml:"mysql" json:"mysql"  mapstructure:"mysql"`
	RedisConfig    `yaml:"redis" json:"redis"  mapstructure:"redis"`
	ConsulConfig   `yaml:"consul" json:"consul"  mapstructure:"consul"`
	StaticSavePath `yaml:"static_save_path" json:"static_save_path" mapstructure:"static_save_path"`
	ServerInfo     `yaml:"server_info" json:"server_info" mapstructure:"server_info"`
}

var (
	TheServerConfig ServerConfig
	TheNacosConfig  NacosConfig
)
