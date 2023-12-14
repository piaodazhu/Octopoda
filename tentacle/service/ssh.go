package service

import (
	"net"

	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/tentacle/config"
	"github.com/piaodazhu/Octopoda/tentacle/logger"
	"github.com/piaodazhu/Octopoda/tentacle/proxy"
)

func SshRegister(conn net.Conn, serialNum uint32, raw []byte) {
	msg := protocols.ProxyMsg{
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
	err = protocols.SendMessageUnique(conn, protocols.TypeSshRegisterResponse, serialNum, payload)
	if err != nil {
		logger.Comm.Println("TypeSshRegisterResponse send error")
	}
}

func SshUnregister(conn net.Conn, serialNum uint32, raw []byte) {
	msg := protocols.ProxyMsg{
		Code: 0,
		Msg:  "OK",
	}
	var payload []byte

	proxy.UnregisterSshService()

	payload, _ = config.Jsoner.Marshal(&msg)
	err := protocols.SendMessageUnique(conn, protocols.TypeSshUnregisterResponse, serialNum, payload)
	if err != nil {
		logger.Comm.Println("TypeSshUnregisterResponse send error")
	}
}
