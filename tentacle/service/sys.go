package service

import (
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"tentacle/config"
	"tentacle/logger"
	"tentacle/message"
	"time"
)

type ScriptParams struct {
	FileName   string
	TargetPath string
	FileBuf    string
	DelayTime  int
}

var shellPath string

func init() {
	shellPath = "/bin/bash"
	_, err := os.Stat(shellPath)
	if err == nil {
		return
	}
	shellPath = "/bin/sh"
	_, err = os.Stat(shellPath)
	if err == nil {
		return
	}
	shellPath = "sh"
}

func RunScript(conn net.Conn, serialNum uint32, raw []byte) {
	sparams := ScriptParams{}
	rmsg := message.Result{
		Rmsg: "OK",
	}
	// var content []byte
	var err error
	// var scriptFile strings.Builder
	// var f *os.File
	var output []byte
	var payload []byte
	var runFunc func()

	if err := config.Jsoner.Unmarshal(raw, &sparams); err != nil {
		logger.Exceptions.Println(err)
		rmsg.Rmsg = "RunScript"
		goto errorout
	}

	runFunc = func() {
		output, err = execScript(&sparams, config.GlobalConfig.Workspace.Root)
		if err != nil {
			rmsg.Rmsg = err.Error()
		} else {
			rmsg.Output = string(output)
		}
	}
	if sparams.DelayTime < 0 {
		runFunc()
	} else {
		time.AfterFunc(time.Duration(sparams.DelayTime)*time.Second, runFunc)
		rmsg.Output = "loaded"
	}

errorout:
	payload, _ = config.Jsoner.Marshal(&rmsg)
	err = message.SendMessageUnique(conn, message.TypeRunScriptResponse, serialNum, payload)
	if err != nil {
		logger.Comm.Println("TypeRunScriptResponse send error")
	}
}

func execScript(sparams *ScriptParams, dir string) ([]byte, error) {
	var content []byte
	var err error
	var scriptFile strings.Builder
	// var f *os.File
	content, err = base64.RawStdEncoding.DecodeString(sparams.FileBuf)
	if err != nil {
		logger.Exceptions.Println("FileDecode")
		return nil, err
	}

	scriptFile.WriteString(config.GlobalConfig.Workspace.Root)
	scriptFile.WriteString(sparams.TargetPath)

	scriptDir := scriptFile.String()
	os.Mkdir(scriptDir, os.ModePerm)
	if scriptFile.String()[scriptFile.Len()-1] != '/' {
		scriptFile.WriteByte('/')
	}

	scriptFile.WriteString(sparams.FileName)
	err = os.WriteFile(scriptFile.String(), content, os.ModePerm)
	if err != nil {
		logger.Exceptions.Println("WriteFile")
		return nil, err
	}
	outputFile := scriptFile.String() + ".output"
	output, _ := os.Create(outputFile)

	// fbuf, _ := os.ReadFile(scriptFile.String())
	// logger.Client.Println(string(fbuf))

	cmd := exec.Command(shellPath, scriptFile.String())
	cmd.Dir = dir
	cmd.Env = append(syscall.Environ(), config.OctopodaEnv(scriptDir, sparams.FileName, outputFile)...)

	scriptErr := cmd.Run()
	if scriptErr != nil {
		logger.Exceptions.Println("Run cmd", err)
	}

	// read output
	// result := []byte{}
	result, err := io.ReadAll(output)
	if err != nil {
		logger.Exceptions.Println(err)
	}
	output.Close()

	// logger.Client.Println("exit code:", cmd.ProcessState.ExitCode())

	os.Remove(scriptFile.String())
	os.Remove(outputFile)

	return result, scriptErr
}

type CommandParams struct {
	Command    string
	Background bool
	DelayTime  int
}

func RunCmd(conn net.Conn, serialNum uint32, raw []byte) {
	rmsg := message.Result{
		Rmsg: "OK",
	}
	var execErr error
	var result []byte = []byte{}
	var payload []byte
	var runFunc func()

	cparams := CommandParams{}
	if err := config.Jsoner.Unmarshal(raw, &cparams); err != nil {
		logger.Exceptions.Println(err)
		rmsg.Rmsg = "RunCmd"
		goto errorout
	}

	runFunc = func() {
		if cparams.Background {
			scriptFile := fmt.Sprintf("%s%d.sh", config.GlobalConfig.Workspace.Root, time.Now().UnixNano())
			content := fmt.Sprintf("(%s) &>/dev/null &", cparams.Command)
			err := os.WriteFile(scriptFile, []byte(content), os.ModePerm)
			if err != nil {
				logger.Exceptions.Println("WriteFile script", err)
			}

			cmd := exec.Command(shellPath, scriptFile)
			cmd.Dir = config.GlobalConfig.Workspace.Root
			cmd.Env = append(syscall.Environ(), config.OctopodaEnv(config.GlobalConfig.Workspace.Root, "NONE", "NONE")...)

			execErr = cmd.Run()
			if execErr != nil {
				logger.Exceptions.Println("Run cmd background", execErr)
			}
			os.Remove(scriptFile)
		} else {
			cmd := exec.Command(shellPath, "-c", cparams.Command)
			cmd.Dir = config.GlobalConfig.Workspace.Root
			cmd.Env = append(syscall.Environ(), config.OctopodaEnv(config.GlobalConfig.Workspace.Root, "NONE", "NONE")...)

			result, execErr = cmd.CombinedOutput()
			if execErr != nil {
				logger.Exceptions.Println("Run cmd foreground", execErr)
			}
		}

		if execErr != nil {
			rmsg.Rmsg = execErr.Error()
		} else {
			rmsg.Output = string(result)
		}
	}

	if cparams.DelayTime < 0 {
		runFunc()
	} else {
		time.AfterFunc(time.Duration(cparams.DelayTime)*time.Second, runFunc)
		rmsg.Output = "loaded"
	}
errorout:
	payload, _ = config.Jsoner.Marshal(&rmsg)
	err := message.SendMessageUnique(conn, message.TypeRunCommandResponse, serialNum, payload)
	if err != nil {
		logger.Comm.Println("TypeRunCommandResponse send error")
	}
}
