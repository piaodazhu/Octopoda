package heartbeat

import (
	"brain/config"
	"brain/logger"
)

type HeartBeatRequest struct {
	Msg string `json:"msg"` // reserved for future usage
	Num uint32 `json:"num"`
}

type HeartBeatResponse struct {
	Msg    string `json:"msg"` // reserved for future usage
	NewNum uint32 `json:"new_num"`
}

func MakeHeartbeat(num uint32) []byte {
	hbInfo := HeartBeatRequest{
		Msg: "ping",
		Num: num,
	}

	serialized_info, err := config.Jsoner.Marshal(hbInfo)
	if err != nil {
		logger.Exceptions.Panic(err)
	}

	return serialized_info
}

func ParseHeartbeat(raw []byte) (HeartBeatRequest, error) {
	info := HeartBeatRequest{}
	err := config.Jsoner.Unmarshal(raw, &info)
	if err != nil {
		logger.Network.Print(err)
		return info, err
	}
	return info, nil
}

func MakeHeartbeatResponse(newNum uint32) []byte {
	hbInfo := HeartBeatResponse{
		Msg: "pong",
		NewNum: newNum,
	}

	serialized_info, err := config.Jsoner.Marshal(hbInfo)
	if err != nil {
		logger.Exceptions.Panic(err)
	}

	return serialized_info
}

func ParseHeartbeatResponse(raw []byte) (HeartBeatResponse, error) {
	rsp := HeartBeatResponse{}
	err := config.Jsoner.Unmarshal(raw, &rsp)
	if err != nil {
		logger.Network.Print(err)
		return rsp, err
	}
	return rsp, nil
}
