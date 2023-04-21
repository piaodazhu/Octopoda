package service

import (
	"encoding/base64"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"tentacle/config"
	"tentacle/logger"
	"tentacle/message"
)

func Reboot() {
	logger.SysInfo.Println("Reboot.")
	_, err := exec.Command("reboot").CombinedOutput()
	if err != nil {
		logger.Exceptions.Fatal(err)
	}
}

func RemoteReboot(conn net.Conn, raw []byte) {
	// valid raw
	err := message.SendMessage(conn, message.TypeCommandResponse, []byte{})
	if err != nil {
		logger.Comm.Println("Reboot send error")
	}

	logger.SysInfo.Println("Reboot.")
	// _, err = exec.Command("reboot").CombinedOutput()
	// if err != nil {
	// 	logger.Client.Fatal(err)
	// }
}

type sshInfo struct {
	Addr     string
	Username string
	Password string
}

func SSHInfo(conn net.Conn, raw []byte) {
	// valid raw

	var addr strings.Builder
	addr.WriteString(config.GlobalConfig.Sshinfo.Ip)
	addr.WriteByte(':')
	addr.WriteString(strconv.Itoa(int(config.GlobalConfig.Sshinfo.Port)))
	sshinfo := sshInfo{
		Addr:     addr.String(),
		Username: config.GlobalConfig.Sshinfo.Username,
		Password: config.GlobalConfig.Sshinfo.Password,
	}
	payload, _ := config.Jsoner.Marshal(&sshinfo)
	err := message.SendMessage(conn, message.TypeCommandResponse, payload)

	if err != nil {
		logger.Comm.Println("SSHInfo send error")
	}
}

type ScriptParams struct {
	FileName   string
	TargetPath string
	FileBuf    string
}

func RunScript(conn net.Conn, raw []byte) {
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

	if err := config.Jsoner.Unmarshal(raw, &sparams); err != nil {
		logger.Exceptions.Println(err)
		rmsg.Rmsg = "FilePush"
		goto errorout
	}

	output, err = execScript(&sparams, config.GlobalConfig.Workspace.Root)
	if err != nil {
		rmsg.Rmsg = err.Error()
	} else {
		rmsg.Rmsg = string(output)
	}

errorout:
	payload, _ = config.Jsoner.Marshal(&rmsg)
	err = message.SendMessage(conn, message.TypeCommandResponse, payload)
	if err != nil {
		logger.Comm.Println("TypeCommandResponse send error")
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

	os.Mkdir(scriptFile.String(), os.ModePerm)
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

	cmd := exec.Command("bash", scriptFile.String())
	cmd.Dir = dir
	cmd.Env = append(syscall.Environ(), config.OctopodaEnv(sparams.TargetPath, sparams.FileName, outputFile)...)

	err = cmd.Run()
	if err != nil {
		logger.Exceptions.Println("Run cmd", err)
		return nil, err
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

	return result, nil
}

func RunCmd(conn net.Conn, raw []byte) {
	rmsg := message.Result{
		Rmsg: "OK",
	}
	var command string
	var err error
	var output []byte
	var payload []byte

	command = string(raw)

	output, err = exec.Command("bash", "-c", command).CombinedOutput()
	if err != nil {
		rmsg.Rmsg = err.Error()
	} else {
		rmsg.Output = string(output)
	}

	payload, _ = config.Jsoner.Marshal(&rmsg)
	err = message.SendMessage(conn, message.TypeCommandResponse, payload)
	if err != nil {
		logger.Comm.Println("TypeCommandResponse send error")
	}
}
