package service

import (
	"net"

	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/tentacle/logger"
)

func InitService() {
	initNodeStatus()
}

func HandleMessage(conn net.Conn) error {
	mtype, serialNum, raw, err := protocols.RecvMessageUnique(conn)
	if err != nil {
		logger.Comm.Println(err)
		return err
	}

	go func() {
		if mtype != protocols.TypeAppLatestVersion {
			logger.Comm.Print(">> Receive Command: ", protocols.MsgTypeString[mtype])
		}

		// If connect is broken in following process, then when we entering
		// this function next time, protocols.RecvMessageUnique(conn) will return error.
		switch mtype {
		case protocols.TypeNodeStatus:
			NodeStatus(conn, serialNum, raw)
		case protocols.TypeFilePush:
			FilePush(conn, serialNum, raw)
		case protocols.TypeFilePull:
			FilePull(conn, serialNum, raw)
		case protocols.TypeFileTree:
			FileTree(conn, serialNum, raw)

		case protocols.TypeNodeLog:
			NodeLog(conn, serialNum, raw)
		case protocols.TypeRunCommand:
			RunCmd(conn, serialNum, raw)
		case protocols.TypeRunScript:
			RunScript(conn, serialNum, raw)

		case protocols.TypeAppCreate:
			AppCreate(conn, serialNum, raw)
		case protocols.TypeAppDelete:
			AppDelete(conn, serialNum, raw)
		case protocols.TypeAppDeploy:
			AppDeploy(conn, serialNum, raw)
		case protocols.TypeAppVersion:
			AppVersions(conn, serialNum, raw)
		case protocols.TypeAppsInfo:
			AppsInfo(conn, serialNum, raw)
		case protocols.TypeAppReset:
			AppReset(conn, serialNum, raw)
		case protocols.TypeAppLatestVersion:
			AppLatestVersion(conn, serialNum, raw)
		case protocols.TypePakmaCommand:
			PakmaCommand(conn, serialNum, raw)
		case protocols.TypeSshRegister:
			SshRegister(conn, serialNum, raw)
		case protocols.TypeSshUnregister:
			SshUnregister(conn, serialNum, raw)
		case protocols.TypeWaitTask:
			TaskWaitResult(conn, serialNum, raw)
		case protocols.TypeCancelTask:
			TaskCancel(conn, serialNum, raw)
		case protocols.TypeQueryTask:
			TaskQueryState(conn, serialNum, raw)
		case protocols.TypeListTasks:
			TaskListAll(conn, serialNum, raw)
		case protocols.TypeAppCommit:
			AppCommit(conn, serialNum, raw)
		default:
			protocols.SendMessageUnique(conn, protocols.TypeUndefined, serialNum, []byte{})
			logger.Comm.Println("unsupported protocol")
		}
	}()

	return nil
}
