package sdk

import (
	"errors"
	"fmt"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/file"
	"github.com/piaodazhu/Octopoda/octl/nameclient"
	"github.com/piaodazhu/Octopoda/octl/node"
	"github.com/piaodazhu/Octopoda/octl/shell"
	"github.com/piaodazhu/Octopoda/protocols"
)

var initalized bool = false

func Init(conf string) error {
	if err := config.InitConfig(conf); err != nil {
		return err
	}

	if err := nameclient.InitClient(); err != nil {
		return err
	}

	initalized = true
	return nil
}

func NodeInfo(name string) (result *protocols.NodeInfo, err error) {
	if !initalized {
		err = errors.New("SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = panicErr.(error)
		}
	}()
	result, err = node.NodeInfo(name)
	return
}

func NodesInfo(names []string) (result *protocols.NodesInfo, err error) {
	if !initalized {
		err = errors.New("SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = panicErr.(error)
		}
	}()
	result, err = node.NodesInfo(names)
	return
}

func NodeStatus(name string) (result *protocols.Status, err error) {
	if !initalized {
		err = errors.New("SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = panicErr.(error)
		}
	}()
	result, err = node.NodeStatus(name)
	return
}

func NodesStatus(names []string) (result *protocols.NodesStatus, err error) {
	if !initalized {
		err = errors.New("SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = panicErr.(error)
		}
	}()
	result, err = node.NodesStatus(names)
	return
}

func DistribFile(localFileOrDir string, targetPath string, names []string) (results []protocols.ExecutionResults, err error) {
	if !initalized {
		err = errors.New("SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = panicErr.(error)
		}
	}()
	results, err = file.DistribFile(localFileOrDir, targetPath, names)
	return
}

func PullFile(pathtype string, node string, fileOrDir string, targetdir string) (result *protocols.ExecutionResults, err error) {
	if !initalized {
		err = errors.New("SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = panicErr.(error)
		}
	}()
	result, err = file.PullFile(pathtype, node, fileOrDir, targetdir)
	return
}

func Run(runstask string, names []string) (results []protocols.ExecutionResults, err error) {
	if !initalized {
		err = errors.New("SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = panicErr.(error)
		}
	}()
	results, err = shell.Run(runstask, names)
	return
}

func XRun(runstask string, names []string, delay int) (results []protocols.ExecutionResults, err error) {
	if !initalized {
		err = errors.New("SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = panicErr.(error)
		}
	}()
	if delay > 0 {
		names = append(names, fmt.Sprintf("-d%d", delay))
	}
	results, err = shell.XRun(runstask, names)
	return
}

func GroupGetAll() (result []string, err error) {
	if !initalized {
		err = errors.New("SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = panicErr.(error)
		}
	}()
	result, err = node.GroupGetAll()
	return
}

func GroupGet(name string) (result *protocols.GroupInfo, err error) {
	if !initalized {
		err = errors.New("SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = panicErr.(error)
		}
	}()
	result, err = node.GroupGet(name)
	return
}

func GroupSet(name string, nocheck bool, names []string) (err error) {
	if !initalized {
		err = errors.New("SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = panicErr.(error)
		}
	}()
	err = node.GroupSet(name, nocheck, names)
	return
}

func GroupDel(name string) (err error) {
	if !initalized {
		err = errors.New("SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = panicErr.(error)
		}
	}()
	err = node.GroupDel(name)
	return
}
