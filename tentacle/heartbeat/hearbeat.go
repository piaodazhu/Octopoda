package heartbeat

import (
	"encoding/json"
	"nworkerd/logger"
	"time"
)

type HeartbeatMessage struct {
	Type int    `json:"type"`
	Raw  string `json:"raw"`
}

type HeartBeatInfo struct {
	Ts  int64 `json:"ts"`
	Cnt int64 `json:"cnt"`
}

type HeartBeatResponse struct {
	Ts  int64 `json:"ts"`
	Cnt int64 `json:"cnt"`
}

func MakeHeartbeat(version int64) []byte {
	hbInfo := HeartBeatInfo{
		Ts:  time.Now().Unix(),
		Cnt: version,
	}

	serialized_info, err := json.Marshal(hbInfo)
	if err != nil {
		logger.Client.Panic(err)
	}

	return serialized_info
}

func ParseHeartbeat(raw []byte) (HeartBeatInfo, error) {
	info := HeartBeatInfo{}
	err := json.Unmarshal(raw, &info)
	if err != nil {
		logger.Client.Print(err)
		return info, nil
	}
	return info, nil
}

func MakeHeartbeatResponse(version int64) []byte {
	return MakeHeartbeat(version)
}

func ParseHeartbeatResponse(raw []byte) (HeartBeatResponse, error) {
	t, err := ParseHeartbeat(raw)
	return HeartBeatResponse(t), err
}
