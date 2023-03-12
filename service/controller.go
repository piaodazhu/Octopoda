package service

import (
	"net"
	"nworkerd/logger"
	"nworkerd/message"
)

func InitService() {
	initNodeState()
}

func HandleConn(conn net.Conn) {
	defer conn.Close()
	mtype, raw, err := message.RecvMessage(conn)
	if err != nil {
		logger.Server.Println(err)
		return 
	}
	switch mtype {
	case message.TypeNodeState: NodeState(conn, raw)
	case message.TypeScenarioVersion: ScenarioVersion(conn, raw)
	case message.TypeFilePush: FilePush(conn, raw)
	default:
		logger.Server.Println("unsupported protocol")
		return 
	}
}
