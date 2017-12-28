package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
	commconf "FishChatServer2/common/conf"
	"FishChatServer2/common/xtime"
)

var (
	confPath string
	Conf     *Config
)

type Config struct {
	commconf.CommConf
	MultiHTTP *commconf.MultiHTTP
	RPCClient *RPCClient
	Redis     *Redis
}

type RPCClient struct {
	RegisterClient *commconf.ServiceDiscoveryClient
}

type Redis struct {
	*commconf.Redis
	Expire xtime.Duration
}

func init() {
	flag.StringVar(&confPath, "conf", "./user-api.toml", "config path")
}

func Init() error {
	if _, err := toml.DecodeFile(confPath, &Conf); err != nil {
		return err
	}
	return nil
}
