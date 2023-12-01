package workgroup

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/viper"
)

var wg workgroupClient
var conf *viper.Viper

func InitWorkgroup(configPath string, client *http.Client) error {
	conf = viper.New()
	conf.SetConfigFile(configPath)
	if err := conf.ReadInConfig(); err != nil {
		fmt.Println(configPath)
		return errors.New("cannot read config: " + err.Error())
	}
	root := conf.GetString("root")
	current := conf.GetString("current")
	password := conf.GetString("password")

	wg = newWorkgroupClient(root, password, current, client)

	if err := wg.auth(); err != nil {
		return errors.New("cannot auth rootgroup: " + err.Error())
	}

	if !wg.valid() {
		return errors.New("current path is invalid")
	}

	return nil
}

func SetHeader(req *http.Request) {
	if req == nil {
		return
	}
	wg.setHeader(req)
}

func Cd(path string) error {
	if err := wg.cd(path); err != nil {
		return err
	}
	conf.Set("current", wg.pwd())
	if err := conf.WriteConfig(); err != nil {
		return err
	}
	return nil
}

func Ls(path string) ([]string, error) {
	return wg.ls(path)
}

func Pwd() string {
	return wg.pwd()
}

func Get(path string) ([]string, error) {
	return wg.get(path)
}

func AddMembers(path string, names ...string) error {
	return wg.addMembers(path, names...)
}

func RemoveMembers(path string, names ...string) error {
	return wg.removeMembers(path, names...)
}

func Grant(path string, password string) error {
	return wg.grant(path, password)
}
