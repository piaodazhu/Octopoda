package service

import (
	"net"
	"tentacle/logger"
	"tentacle/message"
)

func InitService() {
	initNodeStatus()
}

func HandleMessage(conn net.Conn) error {
	mtype, serialNum, raw, err := message.RecvMessageUnique(conn)
	if err != nil {
		logger.Comm.Println(err)
		return err
	}

	go func() {
		if mtype != message.TypeAppLatestVersion {
			logger.Comm.Print(">> Receive Command: ", message.MsgTypeString[mtype])
		}

		// If connect is broken in following process, then when we entering
		// this function next time, message.RecvMessageUnique(conn) will return error.
		switch mtype {
		case message.TypeNodeStatus:
			NodeStatus(conn, serialNum, raw)
		case message.TypeFilePush:
			FilePush(conn, serialNum, raw)
		case message.TypeFilePull:
			FilePull(conn, serialNum, raw)
		case message.TypeFileTree:
			FileTree(conn, serialNum, raw)

		case message.TypeNodeLog:
			NodeLog(conn, serialNum, raw)
		case message.TypeRunCommand:
			RunCmd(conn, serialNum, raw)
		case message.TypeRunScript:
			RunScript(conn, serialNum, raw)

		case message.TypeAppCreate:
			AppCreate(conn, serialNum, raw)
		case message.TypeAppDelete:
			AppDelete(conn, serialNum, raw)
		case message.TypeAppDeploy:
			AppDeploy(conn, serialNum, raw)
		case message.TypeAppVersion:
			AppVersions(conn, serialNum, raw)
		case message.TypeAppsInfo:
			AppsInfo(conn, serialNum, raw)
		case message.TypeAppReset:
			AppReset(conn, serialNum, raw)
		case message.TypeAppLatestVersion:
			AppLatestVersion(conn, serialNum, raw)
		case message.TypePakmaCommand:
			PakmaCommand(conn, serialNum, raw)
		case message.TypeSshRegister:
			SshRegister(conn, serialNum, raw)
		case message.TypeSshUnregister:
			SshUnregister(conn, serialNum, raw)
		case message.TypeWaitTask:
			TaskWaitResult(conn, serialNum, raw)
		case message.TypeCancelTask:
			TaskCancel(conn, serialNum, raw)
		case message.TypeQueryTask:
			TaskQueryState(conn, serialNum, raw)
		case message.TypeListTasks:
			TaskQueryState(conn, serialNum, raw)
		default:
			message.SendMessageUnique(conn, message.TypeUndefined, serialNum, []byte{})
			logger.Comm.Println("unsupported protocol")
		}
	}()

	return nil
}
