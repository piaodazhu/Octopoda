package service

import (
	"net"
	"tentacle/logger"
	"tentacle/message"
)

func InitService() {
	initNodeStatus()
}

func HandleConn(conn net.Conn) {
	defer conn.Close()
	mtype, raw, err := message.RecvMessage(conn)
	if err != nil {
		logger.Comm.Println(err)
		return
	}
	if mtype != message.TypeAppLatestVersion {
		logger.Comm.Print(">> Receive Command: ", message.MsgTypeString[mtype])
	}
	
	switch mtype {
	case message.TypeNodeStatus:
		NodeStatus(conn, raw)
	case message.TypeFilePush:
		FilePush(conn, raw)
	case message.TypeFileTree:
		FileTree(conn, raw)
	case message.TypeNodeLog:
		NodeLog(conn, raw)
	case message.TypeCommandReboot:
		RemoteReboot(conn, raw)
	case message.TypeCommandSSH:
		SSHInfo(conn, raw)
	case message.TypeCommandRun:
		RunCmd(conn, raw)
	case message.TypeCommandRunScript:
		RunScript(conn, raw)

	case message.TypeAppCreate:
		AppCreate(conn, raw)
	case message.TypeAppDelete:
		AppDelete(conn, raw)
	case message.TypeAppDeploy:
		AppDeploy(conn, raw)
	case message.TypeAppVersion:
		AppVersions(conn, raw)
	case message.TypeAppsInfo:
		AppsInfo(conn, raw)
	case message.TypeAppReset:
		AppReset(conn, raw)
	case message.TypeAppLatestVersion:
		AppLatestVersion(conn, raw)
	default:
		logger.Comm.Println("unsupported protocol")
		return
	}
}
