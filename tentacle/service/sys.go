package service

import (
	"encoding/json"
	"net"
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
