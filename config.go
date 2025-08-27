package gouno

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

var GlobalConfig GoUnoConfig

var defaultConfig = "./config/config.yaml"

type GoUnoConfig struct {
	WebServerConfig WebServerConfig `mapstructure:"web_server"`
}

type WebServerConfig struct {
	Debug             bool          `mapstructure:"debug"`
	Address           string        `mapstructure:"address"`
	Port              string        `mapstructure:"port"`
	IdleTimeout       time.Duration `mapstructure:"idle_timeout"`
	ReadTimeout       time.Duration `mapstructure:"read_timeout"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
	WriteTimeout      time.Duration `mapstructure:"write_timeout"`
	RequestTimeout    time.Duration `mapstructure:"request_timeout"`
}

func InitConfig(configFile string) (err error) {

	if configFile == "" {
		configFile = defaultConfig
	}

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	if err = viper.ReadInConfig(); err != nil {
		log.Fatalf("read config failed, err: %v", err)
		return
	}

	if err = viper.Unmarshal(&GlobalConfig); err != nil {
		log.Fatalf("unmarshal config failed, err: %v", err)
		return
	}

	return
}
