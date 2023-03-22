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
	sb.WriteString(config.GlobalConfig.Worker.Ip)
	sb.WriteByte(':')
	sb.WriteString(fmt.Sprint(config.GlobalConfig.Worker.Port))
	listener, err := net.Listen("tcp", sb.String())
	if err != nil {
		logger.Server.Panic(err)
	}
	localListener = listener
}

func ListenAndServe() {
	defer localListener.Close()
	logger.Server.Println("Listening on", localListener.Addr())
	for {
		conn, err := localListener.Accept()
		if err != nil {
			logger.Client.Println(err)
		}
		go service.HandleConn(conn)
	}
}
