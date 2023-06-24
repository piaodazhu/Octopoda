package service

import (
	"net"
	"tentacle/logger"
	"tentacle/message"
	"tentacle/snp"
)

func InitService() {
	initNodeStatus()
}

func HandleConn(conn net.Conn) {
	defer conn.Close()
	mtype, raw, err := message.RecvMessageUnique(conn)
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
	case message.TypeFilePull:
		FilePull(conn, raw)
	case message.TypeFileTree:
		FileTree(conn, raw)

	case message.TypeNodeLog:
		NodeLog(conn, raw)
	case message.TypeRunCommand:
		RunCmd(conn, raw)
	case message.TypeRunScript:
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

func HandleMessage(conn net.Conn) error {
	mtype, raw, err := message.RecvMessageUnique(conn)
	if err != nil {
		logger.Comm.Println(err)
		return err
	}
	if mtype != message.TypeAppLatestVersion {
		logger.Comm.Print(">> Receive Command: ", message.MsgTypeString[mtype])
	}

	// If connect is broken in following process, then when we entering
	// this function next time, message.RecvMessageUnique(conn) will return error.
	switch mtype {
	case message.TypeNodeStatus:
		NodeStatus(conn, raw)
	case message.TypeFilePush:
		FilePush(conn, raw)
	case message.TypeFilePull:
		FilePull(conn, raw)
	case message.TypeFileTree:
		FileTree(conn, raw)

	case message.TypeNodeLog:
		NodeLog(conn, raw)
	case message.TypeRunCommand:
		RunCmd(conn, raw)
	case message.TypeRunScript:
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
	case message.TypePakmaCommand:
		PakmaCommand(conn, raw)
	default:
		message.SendMessageUnique(conn, message.TypeUndefined, snp.GenSerial(), []byte{})
		logger.Comm.Println("unsupported protocol")
	}
	return nil
}
