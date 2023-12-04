package sdk

import (
	"context"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/file"
	"github.com/piaodazhu/Octopoda/octl/httpclient"
	"github.com/piaodazhu/Octopoda/octl/node"
	"github.com/piaodazhu/Octopoda/octl/scenario"
	"github.com/piaodazhu/Octopoda/octl/shell"
	"github.com/piaodazhu/Octopoda/octl/workgroup"
	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/errs"
)

var initalized bool = false

func Init(conf string) *errs.OctlError {
	if err := config.InitConfig(conf); err != nil {
		return err
	}

	if err := httpclient.InitClients(); err != nil {
		return err
	}

	if err := workgroup.InitWorkgroup(httpclient.BrainClient); err != nil {
		return err
	}

	initalized = true
	return nil
}

func NodeInfo(names []string) (result *protocols.NodesInfo, err *errs.OctlError) {
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

func NodeStatus(names []string) (result *protocols.NodesStatus, err *errs.OctlError) {
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

func UploadFile(localFileOrDir string, remoteTargetPath string, isForce bool, names []string) (results []protocols.ExecutionResults, err *errs.OctlError) {
	if !initalized {
		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
		}
	}()
	results, err = file.Upload(localFileOrDir, remoteTargetPath, names, isForce)
	return
}

func DownloadFile(remoteFileOrDir string, localTargetPath string, name string) (result *protocols.ExecutionResults, err *errs.OctlError) {
	if !initalized {
		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
		}
	}()
	result, err = file.Download(remoteFileOrDir, localTargetPath, name)
	return
}

func RunCommand(cmd string, needAlign bool, names []string) (results []protocols.ExecutionResults, err *errs.OctlError) {
	return run(cmd, 0, needAlign, -1, names)
}

func RunScript(scriptFile string, needAlign bool, names []string) (results []protocols.ExecutionResults, err *errs.OctlError) {
	return run(scriptFile, 1, needAlign, -1, names)
}

func RunCommandBackground(cmd string, needAlign bool, names []string) (results []protocols.ExecutionResults, err *errs.OctlError) {
	return run(cmd, 2, needAlign, -1, names)
}

func XRunCommand(cmd string, needAlign bool, delay int, names []string) (results []protocols.ExecutionResults, err *errs.OctlError) {
	if delay < 0 {
		delay = 0
	}
	return run(cmd, 0, needAlign, delay, names)
}

func XRunScript(scriptFile string, needAlign bool, delay int, names []string) (results []protocols.ExecutionResults, err *errs.OctlError) {
	if delay < 0 {
		delay = 0
	}
	return run(scriptFile, 1, needAlign, delay, names)
}

func run(runtask string, cmdType int, needAlign bool, delay int, names []string) (results []protocols.ExecutionResults, err *errs.OctlError) {
	if !initalized {
		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
		}
	}()
	if delay < 0 { // run
		switch cmdType {
		case 0:
			results, err = shell.RunCommand(runtask, false, needAlign, names)
		case 1:
			results, err = shell.RunScript(runtask, needAlign, names)
		case 2:
			results, err = shell.RunCommand(runtask, true, needAlign, names)
		}
	} else {
		switch cmdType {
		case 0:
			results, err = shell.XRunCommand(runtask, false, needAlign, delay, names)
		case 1:
			results, err = shell.XRunScript(runtask, needAlign, delay, names)
		}
	}
	return
}

// func GroupGetAll() (result []string, err *errs.OctlError) {
// 	if !initalized {
// 		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
// 		return
// 	}
// 	defer func() {
// 		if panicErr := recover(); panicErr != nil {
// 			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
// 		}
// 	}()
// 	result, err = node.GroupGetAll()
// 	return
// }

// func GroupGet(name string) (result *protocols.GroupInfo, err *errs.OctlError) {
// 	if !initalized {
// 		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
// 		return
// 	}
// 	defer func() {
// 		if panicErr := recover(); panicErr != nil {
// 			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
// 		}
// 	}()
// 	result, err = node.GroupGet(name)
// 	return
// }

// func GroupSet(name string, nocheck bool, names []string) (err *errs.OctlError) {
// 	if !initalized {
// 		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
// 		return
// 	}
// 	defer func() {
// 		if panicErr := recover(); panicErr != nil {
// 			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
// 		}
// 	}()
// 	err = node.GroupSet(name, nocheck, names)
// 	return
// }

// func GroupDel(name string) (err *errs.OctlError) {
// 	if !initalized {
// 		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
// 		return
// 	}
// 	defer func() {
// 		if panicErr := recover(); panicErr != nil {
// 			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
// 		}
// 	}()
// 	err = node.GroupDel(name)
// 	return
// }

// func Prune(names []string) (err *errs.OctlError) {
// 	if !initalized {
// 		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
// 		return
// 	}
// 	defer func() {
// 		if panicErr := recover(); panicErr != nil {
// 			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
// 		}
// 	}()
// 	err = node.NodesPrune(names)
// 	return
// }

func ScenarioInfo(name string) (rawresult_json []byte, err *errs.OctlError) {
	if !initalized {
		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
		}
	}()
	rawresult_json, err = scenario.ScenarioInfo(name)
	return
}

func ScenariosInfo() (rawresult_json [][]byte, err *errs.OctlError) {
	if !initalized {
		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
		}
	}()
	rawresult_json, err = scenario.ScenariosInfo()
	return
}

func ScenarioVersion(name string, offset, limit int) (rawresult_json []byte, err *errs.OctlError) {
	if !initalized {
		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
		}
	}()
	rawresult_json, err = scenario.ScenarioVersion(name, offset, limit)
	return
}

func NodeAppInfo(name, app, scenario string) (rawresult_json []byte, err *errs.OctlError) {
	if !initalized {
		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
		}
	}()
	rawresult_json, err = node.NodeAppInfo(name, app, scenario)
	return
}

func NodeAppsInfo(name string) (rawresult_json [][]byte, err *errs.OctlError) {
	if !initalized {
		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
		}
	}()
	rawresult_json, err = node.NodeAppsInfo(name)
	return
}

func Apply(ctx context.Context, deployment, target, message string) (logList []string, err *errs.OctlError) {
	// should log
	if !initalized {
		err = errs.New(errs.OctlSdkNotInitializedError, "SDK haven't been initalized")
		return
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = errs.New(errs.OctlSdkPanicRecoverError, panicErr.(error).Error())
		}
	}()
	logList, err = scenario.ScenarioApply(ctx, deployment, target, message)
	return
}
