package sdk

import (
	"errors"
	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/file"
	"github.com/piaodazhu/Octopoda/octl/nameclient"
	"github.com/piaodazhu/Octopoda/octl/node"
	"github.com/piaodazhu/Octopoda/octl/shell"
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

func NodeInfo(name string) (result string, err error) {
	if !initalized {
		err = errors.New("SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); err != nil {
			err = panicErr.(error)
		}
	} ()
	result, err = node.NodeInfo(name)
	return
}

func NodesInfo(names []string) (result string, err error) {
	if !initalized {
		err = errors.New("SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); err != nil {
			err = panicErr.(error)
		}
	} ()
	result, err = node.NodesInfo(names)
	return
}

func DistribFile(localFileOrDir string, targetPath string, names []string) (result string, err error) {
	if !initalized {
		err = errors.New("SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); err != nil {
			err = panicErr.(error)
		}
	} ()
	result, err = file.DistribFile(localFileOrDir, targetPath, names)
	return
}

func Run(runstask string, names []string) (result string, err error) {
	if !initalized {
		err = errors.New("SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); err != nil {
			err = panicErr.(error)
		}
	} ()
	result, err = shell.Run(runstask, names)
	return
}

func XRun(runstask string, names []string) (result string, err error) {
	if !initalized {
		err = errors.New("SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); err != nil {
			err = panicErr.(error)
		}
	} ()
	result, err = shell.XRun(runstask, names)
	return
}

func GroupGetAll() (result string, err error) {
	if !initalized {
		err = errors.New("SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); err != nil {
			err = panicErr.(error)
		}
	} ()
	result, err = node.GroupGetAll()
	return
}

func GroupGet(name string) (result string, err error) {
	if !initalized {
		err = errors.New("SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); err != nil {
			err = panicErr.(error)
		}
	} ()
	result, err = node.GroupGet(name)
	return
}

func GroupSet(name string, nocheck bool, names []string) (result string, err error) {
	if !initalized {
		err = errors.New("SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); err != nil {
			err = panicErr.(error)
		}
	} ()
	result, err = node.GroupSet(name, nocheck, names)
	return
}

func GroupDel(name string) (result string, err error) {
	if !initalized {
		err = errors.New("SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); err != nil {
			err = panicErr.(error)
		}
	} ()
	result, err = node.GroupDel(name)
	return
}
