package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var GlobalConfig ConfigModel

func InitConfig(conf string) {
	if conf != "" {
		viper.SetConfigFile(conf)
		err := viper.ReadInConfig()
		if err != nil {
			panic("cannot read config <" + conf + ">: " + err.Error())
		}
	} else {
		viper.SetConfigName("httpns")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("/etc/octopoda/httpns/")
		viper.AddConfigPath("./")
		err := viper.ReadInConfig()
		if err != nil {
			panic("cannot read config because " + err.Error())
		}
	}

	err := viper.Unmarshal(&GlobalConfig)
	if err != nil {
		panic("cannot unmarshal config because " + err.Error())
	}

	// Create path
	finfo, err := os.Stat(GlobalConfig.Logger.Path)
	if err != nil || !finfo.IsDir() {
		fmt.Println(">> Create ", GlobalConfig.Logger.Path)
		os.Remove(GlobalConfig.Logger.Path)
		os.MkdirAll(GlobalConfig.Logger.Path, os.ModePerm)
	}
}
