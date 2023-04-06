package service

import (
	"encoding/base64"
	"encoding/json"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"tentacle/config"
	"tentacle/logger"
	"tentacle/message"
)

func Reboot() {
	logger.Client.Println("Reboot.")
	_, err := exec.Command("reboot").CombinedOutput()
	if err != nil {
		logger.Client.Fatal(err)
	}
}

func RemoteReboot(conn net.Conn, raw []byte) {
	// valid raw
	err := message.SendMessage(conn, message.TypeCommandResponse, []byte{})
	if err != nil {
		logger.Server.Println("Reboot send error")
	}

	logger.Client.Println("Reboot.")
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
	payload, _ := json.Marshal(&sshinfo)
	err := message.SendMessage(conn, message.TypeCommandResponse, payload)

	if err != nil {
		logger.Server.Println("SSHInfo send error")
	}
}

type ScriptParams struct {
	FileName   string
	TargetPath string
	FileBuf    string
}

func RunScript(conn net.Conn, raw []byte) {
	sparams := ScriptParams{}
	rmsg := RMSG{"OK"}
	// var content []byte
	var err error
	// var scriptFile strings.Builder
	// var f *os.File
	var output []byte
	var payload []byte

	if err := json.Unmarshal(raw, &sparams); err != nil {
		logger.Client.Println(err)
		rmsg.Msg = "FilePush"
		goto errorout
	}

	// content, err = base64.RawStdEncoding.DecodeString(sparams.FileBuf)
	// if err != nil {
	// 	logger.Server.Println("FileDecode")
	// 	rmsg.Msg = "FileDecode"
	// 	goto errorout
	// }

	// scriptFile.WriteString(config.GlobalConfig.Workspace.Store)
	// scriptFile.WriteString(sparams.TargetPath)

	// os.Mkdir(scriptFile.String(), os.ModePerm)

	// scriptFile.WriteString(sparams.FileName)
	// f, err = os.Create(scriptFile.String())
	// if err != nil {
	// 	logger.Server.Println("FileCreate")
	// 	rmsg.Msg = "FileCreate"
	// 	goto errorout
	// }
	// f.Write(content)
	// f.Close()

	// output, err = exec.Command("bash", scriptFile.String()).CombinedOutput()
	// if err != nil {
	// 	rmsg.Msg = err.Error()
	// } else {
	// 	rmsg.Msg = string(output)
	// }

	output, err = execScript(&sparams, []string{"OCTOPODA_ROOT=" + config.GlobalConfig.Workspace.Root})
	if err != nil {
		rmsg.Msg = err.Error()
	} else {
		rmsg.Msg = string(output)
	}

errorout:
	payload, _ = json.Marshal(&rmsg)
	err = message.SendMessage(conn, message.TypeCommandResponse, payload)
	if err != nil {
		logger.Server.Println("TypeCommandResponse send error")
	}
}

func execScript(sparams *ScriptParams, env []string) ([]byte, error) {
	var content []byte
	var err error
	var scriptFile strings.Builder
	// var f *os.File
	content, err = base64.RawStdEncoding.DecodeString(sparams.FileBuf)
	if err != nil {
		logger.Server.Println("FileDecode")
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
	// f, err = os.Create(scriptFile.String())
	if err != nil {
		logger.Server.Println("WriteFile")
		return nil, err
	}
	// f.Write(content)
	// f.Close()

	// fmt.Println(scriptFile.String())
	cmd := exec.Command("bash", scriptFile.String())
	cmd.Env = env
	// fmt.Println(cmd.String(), cmd.Env)
	ret, err := cmd.CombinedOutput()
	os.Remove(scriptFile.String())

	return ret, err
}

func RunCmd(conn net.Conn, raw []byte) {
	rmsg := RMSG{"OK"}
	var command string
	var err error
	var output []byte
	var payload []byte

	command = string(raw)

	output, err = exec.Command("bash", "-c", command).CombinedOutput()
	if err != nil {
		rmsg.Msg = err.Error()
	} else {
		rmsg.Msg = string(output)
	}

	payload, _ = json.Marshal(&rmsg)
	err = message.SendMessage(conn, message.TypeCommandResponse, payload)
	if err != nil {
		logger.Server.Println("TypeCommandResponse send error")
	}
}
