package config

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

var GlobalConfig ConfigModel

func InitConfig(conf string) {
	if conf == "" {
		panic("must specific a tentacle.yaml or brain.yaml config file!")
	}
	confBaseName := filepath.Base(conf)
	parts := strings.Split(confBaseName, ".")
	if len(parts) != 2 || parts[1] != "yaml" {
		panic("invalid config file name!")
	}

	viper.SetConfigFile(conf)
	err := viper.ReadInConfig()
	if err != nil {
		panic("cannot read config <" + conf + ">: " + err.Error())
	}

	err = viper.Unmarshal(&GlobalConfig)
	if err != nil {
		panic("cannot unmarshal config because " + err.Error())
	}

	GlobalConfig.AppName = parts[0]
	GlobalConfig.AppOS = runtime.GOOS
	GlobalConfig.AppArch = runtime.GOARCH
}
