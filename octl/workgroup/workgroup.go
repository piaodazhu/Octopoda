package workgroup

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/protocols/errs"
	"github.com/spf13/viper"
)

var wg workgroupClient
var conf *viper.Viper

func InitWorkgroup(client *http.Client) *errs.OctlError {
	curPathFile := config.GlobalConfig.Workgroup.CurrentPathFile
	root := config.GlobalConfig.Workgroup.Root
	password := config.GlobalConfig.Workgroup.Password

	conf = viper.New()
	conf.SetConfigFile(curPathFile)
	if err := conf.ReadInConfig(); err != nil {
		conf.Set("current", root)
		os.MkdirAll(filepath.Dir(curPathFile), os.ModePerm)
		if err := conf.WriteConfig(); err != nil {
			emsg := "cannot create config: " + err.Error()
			return errs.New(errs.OctlWriteConfigError, emsg)
		}
	}

	current := conf.GetString("current")
	wg = newWorkgroupClient(root, password, current, client)
	if !wg.valid() {
		wg.toRoot()
		conf.Set("current", wg.pwd())
		if err := conf.WriteConfig(); err != nil {
			emsg := "cannot switch to root path: " + err.Error()
			return errs.New(errs.OctlWriteConfigError, emsg)
		}
	}

	if err := wg.auth(); err != nil {
		emsg := "cannot auth rootgroup: " + err.Error()
		return errs.New(errs.OctlWorkgroupAuthError, emsg)
	}
	output.PrintInfof("root workgroup=%s, current workgroup path=%s", wg.root(), wg.pwd())

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
	output.PrintInfof("cd to %s", wg.pwd())
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
	output.PrintString(res)
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
	output.PrintInfoln("OK")
	return nil
}

func RemoveMembers(path string, names ...string) error {
	err := wg.removeMembers(path, names...)
	if err != nil {
		output.PrintFatalln(err)
		return err
	}
	output.PrintInfoln("OK")
	return nil
}

func Grant(path string, password string) error {
	err := wg.grant(path, password)
	if err != nil {
		output.PrintFatalln(err)
		return err
	}
	output.PrintInfoln("OK")
	return nil
}
