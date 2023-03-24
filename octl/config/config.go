package config

import "github.com/spf13/viper"

var GlobalConfig ConfigModel

func InitConfig() {
	viper.SetConfigName("octl")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	viper.AddConfigPath("/etc/octopoda/octl/")
	viper.AddConfigPath("/root/platform/")
	err := viper.ReadInConfig()
	if err != nil {
		panic("cannot read config because " + err.Error())
	}

	err = viper.Unmarshal(&GlobalConfig)
	if err != nil {
		panic("cannot unmarshal config because " + err.Error())
	}
}
