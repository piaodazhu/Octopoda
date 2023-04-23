package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var GlobalConfig ConfigModel

var Stdout bool

func InitConfig() {
	flag.BoolVar(&Stdout, "p", false, "print log to stdout, default is false")
	flag.Parse()

	viper.SetConfigName("brain")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	viper.AddConfigPath("/etc/octopoda/")
	// viper.AddConfigPath("/root/platform/")
	err := viper.ReadInConfig()
	if err != nil {
		panic("cannot read config because " + err.Error())
	}

	err = viper.Unmarshal(&GlobalConfig)
	if err != nil {
		panic("cannot unmarshal config because " + err.Error())
	}

	// if Stdout {
	// 	fmt.Printf("%+v\n", GlobalConfig)
	// }

	// JSON iterator
	if GlobalConfig.JsonFast {
		setFastJsoner()
	} else {
		setStdJsoner()
	}

	// Create path
	finfo, err := os.Stat(GlobalConfig.Workspace.Root)
	if err != nil || !finfo.IsDir() {
		fmt.Println(">> Create ", GlobalConfig.Workspace.Root)
		os.Remove(GlobalConfig.Workspace.Root)
		os.MkdirAll(GlobalConfig.Workspace.Root, os.ModePerm)
	}
	finfo, err = os.Stat(GlobalConfig.Logger.Path)
	if err != nil || !finfo.IsDir() {
		fmt.Println(">> Create ", GlobalConfig.Logger.Path)
		os.Remove(GlobalConfig.Logger.Path)
		os.MkdirAll(GlobalConfig.Logger.Path, os.ModePerm)
	}
	finfo, err = os.Stat(GlobalConfig.Workspace.Store)
	if err != nil || !finfo.IsDir() {
		fmt.Println(">> Create ", GlobalConfig.Workspace.Store)
		os.Remove(GlobalConfig.Workspace.Store)
		os.MkdirAll(GlobalConfig.Workspace.Store, os.ModePerm)
	}
}
