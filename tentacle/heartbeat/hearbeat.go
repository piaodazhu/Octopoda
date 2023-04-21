package heartbeat

import (
	"encoding/json"
	"tentacle/logger"
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

func MakeHeartbeat(ticks int64) []byte {
	hbInfo := HeartBeatInfo{
		Ts:  time.Now().Unix(),
		Cnt: ticks,
	}

	serialized_info, err := json.Marshal(hbInfo)
	if err != nil {
		logger.Exceptions.Panic(err)
	}

	return serialized_info
}

func ParseHeartbeat(raw []byte) (HeartBeatInfo, error) {
	info := HeartBeatInfo{}
	err := json.Unmarshal(raw, &info)
	if err != nil {
		logger.Exceptions.Print(err)
		return info, nil
	}
	return info, nil
}

func MakeHeartbeatResponse(ticks int64) []byte {
	return MakeHeartbeat(ticks)
}

func ParseHeartbeatResponse(raw []byte) (HeartBeatResponse, error) {
	t, err := ParseHeartbeat(raw)
	return HeartBeatResponse(t), err
}
