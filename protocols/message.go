package protocols

import (
	"encoding/binary"
	"fmt"
	"net"
	"github.com/piaodazhu/Octopoda/protocols/security"
	"github.com/piaodazhu/Octopoda/protocols/snp"
)

const (
	TypeNodeJoin = iota
	TypeNodeJoinResponse
	TypeHeartbeat
	TypeHeartbeatResponse

	TypeFilePush
	TypeFilePushResponse

	TypeFilePull
	TypeFilePullResponse

	TypeFileTree
	TypeFileTreeResponse

	TypeNodeStatus
	TypeNodeStatusResponse

	TypeNodeLog
	TypeNodeLogResponse

	TypeRunCommand
	TypeRunCommandResponse

	TypeRunScript
	TypeRunScriptResponse

	TypeAppCreate
	TypeAppCreateResponse

	TypeAppDelete
	TypeAppDeleteResponse

	TypeAppDeploy
	TypeAppDeployResponse

	TypeAppVersion
	TypeAppVersionResponse

	TypeAppsInfo
	TypeAppsInfoResponse

	TypeAppReset
	TypeAppResetResponse

	TypeAppLatestVersion
	TypeAppLatestVersionResponse

	TypePakmaCommand
	TypePakmaCommandResponse

	TypeSshRegister
	TypeSshRegisterResponse

	TypeSshUnregister
	TypeSshUnregisterResponse

	TypeWaitTask
	TypeWaitTaskResponse

	TypeCancelTask
	TypeCancelTaskResponse

	TypeQueryTask
	TypeQueryTaskResponse

	TypeListTasks
	TypeListTasksResponse

	TypeUndefined
)

var MsgTypeString map[int]string = map[int]string{
	TypeNodeJoin:          "NodeJoin",
	TypeNodeJoinResponse:  "NodeJoinResponse",
	TypeHeartbeat:         "Heartbeat",
	TypeHeartbeatResponse: "HeartbeatResponse",

	TypeFilePush:         "FilePush",
	TypeFilePushResponse: "FilePushResponse",

	TypeFilePull:         "TypeFilePull",
	TypeFilePullResponse: "FilePullResponse",

	TypeFileTree:         "FileTree",
	TypeFileTreeResponse: "FileTreeResponse",

	TypeNodeStatus:         "NodeStatus",
	TypeNodeStatusResponse: "NodeStatusResponse",

	TypeNodeLog:         "NodeLog",
	TypeNodeLogResponse: "NodeLogResponse",

	TypeRunCommand:         "RunCommand",
	TypeRunCommandResponse: "RunCommandResponse",

	TypeRunScript:         "RunScript",
	TypeRunScriptResponse: "RunScriptResponse",

	TypeAppCreate:         "AppCreate",
	TypeAppCreateResponse: "AppCreateResponse",

	TypeAppDelete:         "AppDelete",
	TypeAppDeleteResponse: "AppDeleteResponse",

	TypeAppDeploy:         "AppDeploy",
	TypeAppDeployResponse: "AppDeployResponse",

	TypeAppVersion:         "AppVersion",
	TypeAppVersionResponse: "AppVersionResponse",

	TypeAppsInfo:         "AppsInfo",
	TypeAppsInfoResponse: "AppsInfoResponse",

	TypeAppReset:         "AppReset",
	TypeAppResetResponse: "AppResetResponse",

	TypeAppLatestVersion:         "AppLatestVersion",
	TypeAppLatestVersionResponse: "AppLatestVersionResponse",

	TypePakmaCommand:         "TypePakmaCommand",
	TypePakmaCommandResponse: "TypePakmaCommandResponse",

	TypeSshRegister:         "TypeSshRegister",
	TypeSshRegisterResponse: "TypeSshRegisterResponse",

	TypeSshUnregister:         "TypeSshUnregister",
	TypeSshUnregisterResponse: "TypeSshUnregisterResponse",

	TypeWaitTask:         "TypeWaitTask",
	TypeWaitTaskResponse: "TypeWaitTaskResponse",

	TypeCancelTask:         "TypeCancelTask",
	TypeCancelTaskResponse: "TypeCancelTaskResponse",

	TypeQueryTask:         "TypeQueryTask",
	TypeQueryTaskResponse: "TypeQueryTaskResponse",

	TypeListTasks:         "TypeListTasks",
	TypeListTasksResponse: "TypeListTasksResponse",
}

func SendMessageUnique(conn net.Conn, mtype int, serial uint32, raw []byte) error {
	msg := make([]byte, len(raw)+4)
	binary.LittleEndian.PutUint32(msg[0:], serial)
	copy(msg[4:], raw)

	msg, TokenSerial, err := security.AesEncrypt(msg)
	if err != nil {
		return err
	}
	Len := len(msg)
	Buf := make([]byte, Len+16)
	binary.LittleEndian.PutUint32(Buf[0:], uint32(mtype))
	binary.LittleEndian.PutUint32(Buf[4:], uint32(Len))
	binary.LittleEndian.PutUint64(Buf[8:], uint64(TokenSerial))
	copy(Buf[16:], msg)

	Offset := 0
	for Offset < Len+16 {
		n, err := conn.Write(Buf[Offset:])
		if err != nil {
			return err
		}

		Offset += n
	}
	return nil
}

func RecvMessageUnique(conn net.Conn) (int, uint32, []byte, error) {
	Len := 0
	Buf := make([]byte, 16)

	Offset := 0
	for Offset < 16 {
		n, err := conn.Read(Buf[Offset:])
		if err != nil {
			return 0, 0, nil, err
		}
		Offset += n
	}

	mtype := int(binary.LittleEndian.Uint32(Buf[0:]))
	Len = int(binary.LittleEndian.Uint32(Buf[4:]))
	TokenSerial := int64(binary.LittleEndian.Uint64(Buf[8:]))

	Buf = make([]byte, Len)
	Offset = 0
	for Offset < Len {
		n, err := conn.Read(Buf[Offset:])
		if err != nil {
			return 0, 0, nil, err
		}

		Offset += n
	}
	Buf, err := security.AesDecrypt(Buf, TokenSerial)
	if err != nil {
		return 0, 0, nil, err
	}
	serial := binary.LittleEndian.Uint32(Buf[0:])
	if !snp.CheckSerial(serial) {
		return 0, 0, nil, fmt.Errorf("duplicated message")
	}
	return mtype, serial, Buf[4:], nil
}