package message

import (
	"brain/logger"
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

	TypeNodeState
	TypeNodeStateResponse

	TypeScenarioVersion
	TypeScenarioVersionResponse

	TypeModVersion
	TypeModVersionResponse
)

func SendMessage(conn net.Conn, mtype int, raw []byte) error {
	Len := len(raw)
	Buf := make([]byte, Len+4)
	binary.LittleEndian.PutUint16(Buf[0:], uint16(mtype))
	binary.LittleEndian.PutUint16(Buf[2:], uint16(Len))
	copy(Buf[4:], raw)

	Offset := 0
	for Offset < Len {
		n, err := conn.Write(Buf[Offset:])
		if err != nil {
			logger.Tentacle.Print(err)
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
		if err != nil {
			logger.Tentacle.Print(err)
			return 0, nil, err
		}

		Offset += n
	}

	mtype := int(binary.LittleEndian.Uint16(Buf[0:]))
	Len = int(binary.LittleEndian.Uint16(Buf[2:]))
	Buf = make([]byte, Len)
	Offset = 0
	for Offset < Len {
		n, err := conn.Read(Buf[Offset:])
		if err != nil {
			logger.Tentacle.Println(err)
			return 0, nil, err
		}

		Offset += n
	}
	return mtype, Buf, nil
}
