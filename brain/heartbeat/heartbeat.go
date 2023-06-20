package heartbeat

import (
	"brain/config"
	"brain/logger"
)

type HeartBeatInfo struct {
	Msg string `json:"msg"` // reserved for future usage
}

type HeartBeatResponse struct {
	Msg string `json:"msg"` // reserved for future usage
}

func MakeHeartbeat(msg string) []byte {
	hbInfo := HeartBeatInfo{
		Msg: msg,
	}

	serialized_info, err := config.Jsoner.Marshal(hbInfo)
	if err != nil {
		logger.Exceptions.Panic(err)
	}

	return serialized_info
}

func ParseHeartbeat(raw []byte) (HeartBeatInfo, error) {
	info := HeartBeatInfo{}
	err := config.Jsoner.Unmarshal(raw, &info)
	if err != nil {
		logger.Network.Print(err)
		return info, nil
	}
	return info, nil
}

func MakeHeartbeatResponse(msg string) []byte {
	return MakeHeartbeat(msg)
}

func ParseHeartbeatResponse(raw []byte) (HeartBeatResponse, error) {
	t, err := ParseHeartbeat(raw)
	return HeartBeatResponse(t), err
}
