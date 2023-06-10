package config

import (
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
		viper.SetConfigName("octl")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("/etc/octopoda/octl/")
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

	if GlobalConfig.JsonFast {
		setFastJsoner()
	} else {
		setStdJsoner()
	}
}
