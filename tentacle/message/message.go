package message

import (
	"encoding/binary"
	"net"
)

// type Message struct {
// 	Type int32
// 	Len  int32
// 	Raw  string
// }

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
	
	TypeCommandReboot
	TypeCommandSSH
	TypeCommandRun
	TypeCommandRunScript
	TypeCommandResponse

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

	TypeUndefined
)

var MsgTypeString map[int]string = map[int]string {
	TypeNodeJoin: "NodeJoin",
	TypeNodeJoinResponse: "NodeJoinResponse",
	TypeHeartbeat: "Heartbeat",
	TypeHeartbeatResponse: "HeartbeatResponse",

	TypeFilePush: "FilePush",
	TypeFilePushResponse: "FilePushResponse",

	TypeFilePull: "TypeFilePull",
	TypeFilePullResponse: "FilePullResponse",

	TypeFileTree: "FileTree",
	TypeFileTreeResponse: "FileTreeResponse",

	TypeNodeStatus: "NodeState",
	TypeNodeStatusResponse: "NodeStateResponse",

	TypeNodeLog: "NodeLog",
	TypeNodeLogResponse: "NodeLogResponse",
	
	TypeCommandReboot: "CommandReboot",
	TypeCommandSSH: "CommandSSH",
	TypeCommandRun: "CommandRun",
	TypeCommandRunScript: "CommandRunScript",
	TypeCommandResponse: "CommandResponse",

	TypeAppCreate: "AppCreate",
	TypeAppCreateResponse: "AppCreateResponse",

	TypeAppDelete: "AppDelete",
	TypeAppDeleteResponse: "AppDeleteResponse",

	TypeAppDeploy: "AppDeploy",
	TypeAppDeployResponse: "AppDeployResponse",

	TypeAppVersion: "AppVersion",
	TypeAppVersionResponse: "AppVersionResponse",

	TypeAppsInfo: "AppsInfo",
	TypeAppsInfoResponse: "AppsInfoResponse",

	TypeAppReset: "AppReset",
	TypeAppResetResponse: "AppResetResponse",

	TypeAppLatestVersion: "AppLatestVersion",
	TypeAppLatestVersionResponse: "AppLatestVersionResponse",
}

func SendMessage(conn net.Conn, mtype int, raw []byte) error {
	Len := len(raw)
	Buf := make([]byte, Len+8)
	binary.LittleEndian.PutUint32(Buf[0:], uint32(mtype))
	binary.LittleEndian.PutUint32(Buf[4:], uint32(Len))
	copy(Buf[8:], raw)

	Offset := 0
	for Offset < Len+8 {
		n, err := conn.Write(Buf[Offset:])
		if err != nil {
			return err
		}

		Offset += n
	}
	return nil
}

func RecvMessage(conn net.Conn) (int, []byte, error) {
	Len := 0
	Buf := make([]byte, 8)

	Offset := 0
	for Offset < 8 {
		n, err := conn.Read(Buf[Offset:])
		if err != nil {
			return 0, nil, err
		}

		Offset += n
	}

	mtype := int(binary.LittleEndian.Uint32(Buf[0:]))
	Len = int(binary.LittleEndian.Uint32(Buf[4:]))
	Buf = make([]byte, Len)
	Offset = 0
	for Offset < Len {
		n, err := conn.Read(Buf[Offset:])
		if err != nil {
			return 0, nil, err
		}

		Offset += n
	}
	return mtype, Buf, nil
}
