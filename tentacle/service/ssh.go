package service

import (
	"net"
	"tentacle/config"
	"tentacle/logger"
	"tentacle/message"
	"tentacle/proxy"
	"tentacle/snp"
)

type proxyMsg struct {
	Code int
	Msg  string
	Data string
}

func SshRegister(conn net.Conn, raw []byte) {
	msg := proxyMsg{
		Code: 0,
		Msg:  "OK",
	}
	var payload []byte

	serviceAddr, err := proxy.RegisterSshService()
	if err != nil {
		msg.Code = 1
		msg.Msg = err.Error()
		goto errorout
	}
	msg.Data = serviceAddr

errorout:
	payload, _ = config.Jsoner.Marshal(&msg)
	err = message.SendMessageUnique(conn, message.TypeSshRegisterResponse, snp.GenSerial(), payload)
	if err != nil {
		logger.Comm.Println("TypeSshRegisterResponse send error")
	}
}

func SshUnregister(conn net.Conn, raw []byte) {
	msg := proxyMsg{
		Code: 0,
		Msg:  "OK",
	}
	var payload []byte

	proxy.UnregisterSshService()

	payload, _ = config.Jsoner.Marshal(&msg)
	err := message.SendMessageUnique(conn, message.TypeSshUnregisterResponse, snp.GenSerial(), payload)
	if err != nil {
		logger.Comm.Println("TypeSshUnregisterResponse send error")
	}
}
