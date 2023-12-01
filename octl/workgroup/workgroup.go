package workgroup

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/piaodazhu/Octopoda/octl/output"
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
		output.PrintFatalln(err)
		return err
	}
	conf.Set("current", wg.pwd())
	if err := conf.WriteConfig(); err != nil {
		output.PrintFatalln(err)
		return err
	}
	return nil
}

func Ls(path string) ([]string, error) {
	children, err := wg.ls(path)
	if err != nil {
		output.PrintFatalln(err)
		return nil, err
	}
	output.PrintJSON(children)
	return children, nil
}

func Pwd() string {
	res := wg.pwd()
	if len(res) == 0 {
		res = "/"
	}
	fmt.Println(res)
	output.PrintJSON(res)
	return res
}

func Get(path string) ([]string, error) {
	members, err := wg.get(path)
	if err != nil {
		output.PrintFatalln(err)
		return nil, err
	}
	output.PrintJSON(members)
	return members, nil
}

func AddMembers(path string, names ...string) error {
	err := wg.addMembers(path, names...)
	if err != nil {
		output.PrintFatalln(err)
		return err
	}
	return nil
}

func RemoveMembers(path string, names ...string) error {
	err := wg.removeMembers(path, names...)
	if err != nil {
		output.PrintFatalln(err)
		return err
	}
	return nil
}

func Grant(path string, password string) error {
	err := wg.grant(path, password)
	if err != nil {
		output.PrintFatalln(err)
		return err
	}
	return nil
}
