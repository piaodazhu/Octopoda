package config

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

var GlobalConfig ConfigModel

var Stdout bool

func InitConfig() {
	flag.BoolVar(&Stdout, "p", false, "print log to stdout, default is false")
	flag.Parse()

	viper.SetConfigName("octopoda")
	viper.SetConfigFile("yaml")
	viper.AddConfigPath("./")
	viper.AddConfigPath("/etc/brain/")
	viper.AddConfigPath("/root/platform/")
	err := viper.ReadInConfig()
	if err != nil {
		panic("cannot read config because " + err.Error())
	}

	err = viper.Unmarshal(&GlobalConfig)
	if err != nil {
		panic("cannot unmarshal config because " + err.Error())
	}

	if Stdout {
		fmt.Printf("%+v\n", GlobalConfig)
	}
}
