package network

import (
	"fmt"
	"net"
	"tentacle/config"
	"tentacle/logger"
	"tentacle/service"
	"strings"
)

var localListener net.Listener

func InitListener() {
	var sb strings.Builder
	sb.WriteString(config.GlobalConfig.Ip)
	sb.WriteByte(':')
	sb.WriteString(fmt.Sprint(config.GlobalConfig.Port))
	listener, err := net.Listen("tcp", sb.String())
	if err != nil {
		logger.Exceptions.Panic(err)
	}
	localListener = listener
}

func ListenAndServe() {
	defer localListener.Close()
	logger.SysInfo.Println("Listening on", localListener.Addr())
	for {
		conn, err := localListener.Accept()
		if err != nil {
			logger.Exceptions.Println(err)
		}
		go service.HandleConn(conn)
	}
}
