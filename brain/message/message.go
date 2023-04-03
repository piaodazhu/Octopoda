package message

import (
	"encoding/binary"
	"net"
)

// type Message struct {
// 	Type int16
// 	Len  int16
// 	Raw  string
// }

const (
	TypeNodeJoin = iota
	TypeNodeJoinResponse
	TypeHeartbeat
	TypeHeartbeatResponse

	TypeFilePush
	TypeFilePushResponse

	TypeFileTree
	TypeFileTreeResponse

	TypeNodeState
	TypeNodeStateResponse

	TypeScenarioVersion
	TypeScenarioVersionResponse

	TypeModVersion
	TypeModVersionResponse

	TypeNodeLog
	TypeNodeLogResponse

	TypeAppCreate
	TypeAppCreateResponse

	TypeAppDelete
	TypeAppDeleteResponse

	TypeCommandReboot
	TypeCommandSSH
	TypeCommandRun
	TypeCommandRunScript
	TypeCommandResponse

	TypeAppDeploy
	TypeAppDeployResponse

	TypeAppVersion
	TypeAppVersionResponse

	TypeAppsInfo
	TypeAppsInfoResponse

	TypeAppReset
	TypeAppResetResponse
)

func SendMessage(conn net.Conn, mtype int, raw []byte) error {
	Len := len(raw)
	Buf := make([]byte, Len+4)
	binary.LittleEndian.PutUint16(Buf[0:], uint16(mtype))
	binary.LittleEndian.PutUint16(Buf[2:], uint16(Len))
	copy(Buf[4:], raw)

	Offset := 0
	for Offset < Len+4 {
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
	Buf := make([]byte, 4)

	Offset := 0
	for Offset < 4 {
		n, err := conn.Read(Buf[Offset:])
		// logger.Tentacle.Print(n, err)
		if err != nil {
			return 0, nil, err
		}

		Offset += n
	}

	mtype := int(binary.LittleEndian.Uint16(Buf[0:]))
	Len = int(binary.LittleEndian.Uint16(Buf[2:]))
	Buf = make([]byte, Len)
	Offset = 0
	// logger.Tentacle.Print(Len, mtype)
	for Offset < Len {
		n, err := conn.Read(Buf[Offset:])
		if err != nil {
			return 0, nil, err
		}

		Offset += n
	}
	// logger.Tentacle.Print(Buf)
	return mtype, Buf, nil
}
