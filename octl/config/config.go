package config

import (
	"github.com/piaodazhu/Octopoda/protocols/errs"
	"github.com/spf13/viper"
)

var GlobalConfig ConfigModel

func InitConfig(conf string) *errs.OctlError {
	if conf != "" {
		viper.SetConfigFile(conf)
		err := viper.ReadInConfig()
		if err != nil {
			return errs.New(errs.OctlReadConfigError, "cannot read config <"+conf+">: "+err.Error())
		}
	} else {
		viper.SetConfigName("octl")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("/etc/octopoda/octl/")
		viper.AddConfigPath("./")
		err := viper.ReadInConfig()
		if err != nil {
			return errs.New(errs.OctlReadConfigError, "cannot read config because "+err.Error())
		}
	}

	err := viper.Unmarshal(&GlobalConfig)
	if err != nil {
		return errs.New(errs.OctlReadConfigError, "cannot unmarshal config because "+err.Error())
	}

	if GlobalConfig.JsonFast {
		setFastJsoner()
	} else {
		setStdJsoner()
	}
	return nil
}
