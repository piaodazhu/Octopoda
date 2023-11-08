package sdk

import (
	"fmt"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/file"
	"github.com/piaodazhu/Octopoda/octl/nameclient"
	"github.com/piaodazhu/Octopoda/octl/node"
	"github.com/piaodazhu/Octopoda/octl/shell"
	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/errs"
)

var initalized bool = false

func Init(conf string) *errs.OctlError {
	if err := config.InitConfig(conf); err != nil {
		return err
	}

	if err := nameclient.InitClient(); err != nil {
		return err
	}

	initalized = true
	return nil
}

func NodeInfo(name string) (result *protocols.NodeInfo, err *errs.OctlError) {
	if !initalized {
		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
		}
	}()
	result, err = node.NodeInfo(name)
	return
}

func NodesInfo(names []string) (result *protocols.NodesInfo, err *errs.OctlError) {
	if !initalized {
		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
		}
	}()
	result, err = node.NodesInfo(names)
	return
}

func NodeStatus(name string) (result *protocols.Status, err *errs.OctlError) {
	if !initalized {
		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
		}
	}()
	result, err = node.NodeStatus(name)
	return
}

func NodesStatus(names []string) (result *protocols.NodesStatus, err *errs.OctlError) {
	if !initalized {
		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
		}
	}()
	result, err = node.NodesStatus(names)
	return
}

func DistribFile(localFileOrDir string, targetPath string, names []string) (results []protocols.ExecutionResults, err *errs.OctlError) {
	if !initalized {
		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
		}
	}()
	results, err = file.DistribFile(localFileOrDir, targetPath, names)
	return
}

func PullFile(pathtype string, node string, fileOrDir string, targetdir string) (result *protocols.ExecutionResults, err *errs.OctlError) {
	if !initalized {
		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
		}
	}()
	result, err = file.PullFile(pathtype, node, fileOrDir, targetdir)
	return
}

func Run(runstask string, names []string) (results []protocols.ExecutionResults, err *errs.OctlError) {
	if !initalized {
		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
		}
	}()
	results, err = shell.Run(runstask, names)
	return
}

func XRun(runstask string, names []string, delay int) (results []protocols.ExecutionResults, err *errs.OctlError) {
	if !initalized {
		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
		}
	}()
	if delay > 0 {
		names = append(names, fmt.Sprintf("-d%d", delay))
	}
	results, err = shell.XRun(runstask, names)
	return
}

func GroupGetAll() (result []string, err *errs.OctlError) {
	if !initalized {
		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
		}
	}()
	result, err = node.GroupGetAll()
	return
}

func GroupGet(name string) (result *protocols.GroupInfo, err *errs.OctlError) {
	if !initalized {
		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
		}
	}()
	result, err = node.GroupGet(name)
	return
}

func GroupSet(name string, nocheck bool, names []string) (err *errs.OctlError) {
	if !initalized {
		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
		}
	}()
	err = node.GroupSet(name, nocheck, names)
	return
}

func GroupDel(name string) (err *errs.OctlError) {
	if !initalized {
		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
		}
	}()
	err = node.GroupDel(name)
	return
}
