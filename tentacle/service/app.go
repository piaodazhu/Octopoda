package service

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/san"
	"github.com/piaodazhu/Octopoda/tentacle/app"
	"github.com/piaodazhu/Octopoda/tentacle/config"
	"github.com/piaodazhu/Octopoda/tentacle/logger"
	"github.com/piaodazhu/Octopoda/tentacle/task"
)

func AppsInfo(conn net.Conn, serialNum uint32, raw []byte) {
	err := protocols.SendMessageUnique(conn, protocols.TypeAppsInfoResponse, serialNum, app.Digest())
	if err != nil {
		logger.Comm.Println("AppsInfo send error")
	}
}

func AppInfo(conn net.Conn, serialNum uint32, raw []byte) {
	aParams := san.AppBasic{}
	err := config.Jsoner.Unmarshal(raw, &aParams)
	var payload []byte = []byte("[]")
	if err != nil {
		logger.Exceptions.Println(err)
		goto errorout
	}
	payload = app.AppInfo(aParams.Name, aParams.Scenario)
errorout:
	err = protocols.SendMessageUnique(conn, protocols.TypeAppInfoResponse, serialNum, payload)
	if err != nil {
		logger.Comm.Println("AppsInfo send error")
	}
}

func AppCreate(conn net.Conn, serialNum uint32, raw []byte) {
	acParams := san.AppCreateParams{}
	if err := config.Jsoner.Unmarshal(raw, &acParams); err != nil {
		logger.Exceptions.Println("invalid arguments: ", err)
		// SNED BACK
		err = protocols.SendMessageUnique(conn, protocols.TypeAppCreateResponse, serialNum, []byte{})
		if err != nil {
			logger.Comm.Println("TypeAppCreateResponse send error")
		}
		return
	}

	var utaskFunc func() *protocols.Result
	var ucancelFunc func()

	utaskFunc = func() *protocols.Result {
		rmsg := protocols.Result{
			Rmsg: "OK",
		}
		rmsg.Modified = false
		fullname := acParams.Name + "@" + acParams.Scenario
		var version san.Version

		// is exists?
		if !app.Exists(acParams.Name, acParams.Scenario) {
			if !app.Create(acParams.Name, acParams.Scenario, acParams.Description) {
				logger.Exceptions.Println("app.Create")
				rmsg.Rmsg = "Failed Create App"
				return &rmsg
			}
			if !app.GitCreate(fullname) {
				logger.Exceptions.Println("app.GitCreate")
				rmsg.Rmsg = "Failed Init the Repo"
				return &rmsg
			}
		}
		// unpack files
		err := unpackFilesNoWrap(acParams.FilePack, config.GlobalConfig.Workspace.Root+fullname)
		if err != nil {
			logger.Exceptions.Println("unpack Files")
			rmsg.Rmsg = err.Error()
			return &rmsg
		}
		// commit
		version, err = app.GitCommit(fullname, acParams.Message)
		if err != nil {
			logger.Exceptions.Println("app.GitCommit")
			if _, ok := err.(app.EmptyCommitError); ok {
				rmsg.Rmsg = "OK: No Change"
				rmsg.Version = app.CurVersion(acParams.Name, acParams.Scenario).Hash
			} else {
				rmsg.Rmsg = err.Error()
			}
			return &rmsg
		}

		// update nodeApps
		if !app.Update(acParams.Name, acParams.Scenario, version) {
			logger.Exceptions.Println("app.Update")
			rmsg.Rmsg = "Faild to update app version"
		}
		rmsg.Rmsg = "OK"
		rmsg.Version = version.Hash
		rmsg.Modified = true
		app.Save()
		return &rmsg
	}

	ucancelFunc = func() {
		// 不好控制回滚，状态太复杂
		// 肯定不会阻塞，所以这里只是一个假Kill
	}

	brief := fmt.Sprintf("%s@%s(%s)", acParams.Name, acParams.Scenario, acParams.Description)
	taskId, err := task.TaskManager.CreateTask(brief, utaskFunc, ucancelFunc)
	if err != nil {
		// ERROR
		logger.Exceptions.Println("cannot create task: ", err)
		err = protocols.SendMessageUnique(conn, protocols.TypeAppCreateResponse, serialNum, []byte{})
		if err != nil {
			logger.Comm.Println("TypeAppCreateResponse send error")
		}
		return
	}
	err = protocols.SendMessageUnique(conn, protocols.TypeAppCreateResponse, serialNum, []byte(taskId))
	if err != nil {
		logger.Comm.Println("TypeAppCreateResponse send error")
	}
}

func AppDeploy(conn net.Conn, serialNum uint32, raw []byte) {
	adParams := san.AppDeployParams{}
	if err := config.Jsoner.Unmarshal(raw, &adParams); err != nil {
		logger.Exceptions.Println("invalid arguments: ", err)
		// SNED BACK
		err = protocols.SendMessageUnique(conn, protocols.TypeAppDeployResponse, serialNum, []byte{})
		if err != nil {
			logger.Comm.Println("TypeAppDeployResponse send error")
		}
		return
	}

	var utaskFunc func() *protocols.Result
	var ucancelFunc func()
	cmdChan := make(chan *exec.Cmd, 1)

	utaskFunc = func() *protocols.Result {
		rmsg := protocols.Result{
			Rmsg: "OK",
		}

		rmsg.Modified = false
		fullname := adParams.Name + "@" + adParams.Scenario
		var version san.Version

		// is new?
		if !app.Exists(adParams.Name, adParams.Scenario) {
			if !app.Create(adParams.Name, adParams.Scenario, adParams.Description) {
				logger.Exceptions.Println("app.Create")
				rmsg.Rmsg = "Failed Create App"
				return &rmsg
			}
			if !app.GitCreate(fullname) {
				logger.Exceptions.Println("app.GitCreate")
				rmsg.Rmsg = "Failed Init the Repo"
				return &rmsg
			}
		}
		// run script
		sparams := protocols.ScriptParams{
			FileName:   fmt.Sprintf("script_%s.sh", time.Now().Format("2006_01_02_15_04")),
			TargetPath: "scripts/",
			FileBuf:    adParams.Script,
		}
		output, err := execScript(&sparams, config.GlobalConfig.Workspace.Root+fullname, cmdChan)
		rmsg.Output = string(output)

		if err != nil {
			logger.Exceptions.Println("execScript")
			rmsg.Rmsg = err.Error()
			return &rmsg
		}

		version, isClean, err := app.GitStatus(fullname)
		if err != nil {
			rmsg.Rmsg = err.Error()
			rmsg.Version = app.CurVersion(adParams.Name, adParams.Scenario).Hash
			return &rmsg
		}
		rmsg.Version = version.Hash
		rmsg.Modified = !isClean

		return &rmsg
	}

	ucancelFunc = func() {
		cmd := <-cmdChan
		cmd.Process.Kill()
	}

	brief := fmt.Sprintf("%s@%s(%s)", adParams.Name, adParams.Scenario, adParams.Description)
	taskId, err := task.TaskManager.CreateTask(brief, utaskFunc, ucancelFunc)
	if err != nil {
		// ERROR
		logger.Exceptions.Println("cannot create task: ", err)
		err = protocols.SendMessageUnique(conn, protocols.TypeAppDeployResponse, serialNum, []byte{})
		if err != nil {
			logger.Comm.Println("TypeAppDeployResponse send error")
		}
		return
	}
	err = protocols.SendMessageUnique(conn, protocols.TypeAppDeployResponse, serialNum, []byte(taskId))
	if err != nil {
		logger.Comm.Println("TypeAppDeployResponse send error")
	}
}

func AppCommit(conn net.Conn, serialNum uint32, raw []byte) {
	acParams := san.AppBasic{}
	if err := config.Jsoner.Unmarshal(raw, &acParams); err != nil {
		logger.Exceptions.Println("invalid arguments: ", err)
		// SNED BACK
		err = protocols.SendMessageUnique(conn, protocols.TypeAppCommitResponse, serialNum, []byte{})
		if err != nil {
			logger.Comm.Println("TypeAppCommitResponse send error")
		}
		return
	}

	var utaskFunc func() *protocols.Result
	var ucancelFunc func()
	cmdChan := make(chan *exec.Cmd, 1)

	utaskFunc = func() *protocols.Result {
		rmsg := protocols.Result{
			Rmsg: "OK",
		}

		rmsg.Modified = false
		fullname := acParams.Name + "@" + acParams.Scenario
		var version san.Version

		// is new?
		if !app.Exists(acParams.Name, acParams.Scenario) {
			logger.Exceptions.Println("app.Create")
			rmsg.Rmsg = "App repo not exist"
			return &rmsg
		}

		// append a dummy file
		err := appendDummyFile(fullname)
		if err != nil {
			logger.Exceptions.Println("append dummy file", err)
		}

		// commit
		version, err = app.GitCommit(fullname, acParams.Message)
		if err != nil {
			if _, ok := err.(app.EmptyCommitError); ok {
				rmsg.Rmsg = "OK: No Change"
				rmsg.Version = app.CurVersion(acParams.Name, acParams.Scenario).Hash
			} else {
				rmsg.Rmsg = err.Error()
				logger.Exceptions.Println("app.GitCommit")
			}
			return &rmsg
		}

		// update nodeApps
		if !app.Update(acParams.Name, acParams.Scenario, version) {
			logger.Exceptions.Println("app.Update")
			rmsg.Rmsg = "Faild to update app version"
		}

		rmsg.Version = version.Hash
		rmsg.Modified = true
		app.Save()

		return &rmsg
	}

	ucancelFunc = func() {
		cmd := <-cmdChan
		cmd.Process.Kill()
	}

	brief := fmt.Sprintf("%s@%s(%s)", acParams.Name, acParams.Scenario, acParams.Description)
	taskId, err := task.TaskManager.CreateTask(brief, utaskFunc, ucancelFunc)
	if err != nil {
		// ERROR
		logger.Exceptions.Println("cannot create task: ", err)
		err = protocols.SendMessageUnique(conn, protocols.TypeAppCommitResponse, serialNum, []byte{})
		if err != nil {
			logger.Comm.Println("TypeAppCommitResponse send error")
		}
		return
	}
	err = protocols.SendMessageUnique(conn, protocols.TypeAppCommitResponse, serialNum, []byte(taskId))
	if err != nil {
		logger.Comm.Println("TypeAppCommitResponse send error")
	}
}

type AppDeleteParams struct {
	Name     string
	Scenario string
}

func AppDelete(conn net.Conn, serialNum uint32, raw []byte) {
	adParams := &AppDeleteParams{}
	rmsg := protocols.Result{
		Rmsg: "OK",
	}
	var payload []byte
	var ok bool
	var subdirName string

	err := config.Jsoner.Unmarshal(raw, adParams)
	if err != nil {
		logger.Exceptions.Println(err)
		rmsg.Rmsg = "Invalid Params"
		goto errorout
	}
	subdirName = adParams.Name + "@" + adParams.Scenario

	// is new?
	if ok = app.Exists(adParams.Name, adParams.Scenario); ok {
		if !app.Delete(adParams.Name, adParams.Scenario) {
			logger.Exceptions.Println("app.Delete")
			rmsg.Rmsg = "Failed Delete App"
			goto errorout
		}
	}
	os.RemoveAll(config.GlobalConfig.Workspace.Root + subdirName)
	app.Save()
errorout:
	payload, _ = config.Jsoner.Marshal(&rmsg)
	err = protocols.SendMessageUnique(conn, protocols.TypeAppDeleteResponse, serialNum, payload)
	if err != nil {
		logger.Comm.Println("AppDelete send error")
	}
}

func appendDummyFile(fullname string) error {
	f, err := os.OpenFile(config.GlobalConfig.Workspace.Root+fullname+"/.DUMMY", os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModeAppend)
	if err != nil {
		return err
	}
	_, err = f.WriteString("+\n")
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}
