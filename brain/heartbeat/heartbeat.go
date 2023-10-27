package heartbeat

import (
	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/logger"
	"github.com/piaodazhu/Octopoda/protocols"
)

func MakeHeartbeat(num uint32) []byte {
	hbInfo := protocols.HeartBeatRequest{
		Msg: "ping",
		Num: num,
	}

	serialized_info, err := config.Jsoner.Marshal(hbInfo)
	if err != nil {
		logger.Exceptions.Panic(err)
	}

	return serialized_info
}

func ParseHeartbeat(raw []byte) (protocols.HeartBeatRequest, error) {
	info := protocols.HeartBeatRequest{}
	err := config.Jsoner.Unmarshal(raw, &info)
	if err != nil {
		logger.Network.Print(err)
		return info, err
	}
	return info, nil
}

func MakeHeartbeatResponse(newNum uint32) []byte {
	hbInfo := protocols.HeartBeatResponse{
		Msg:    "pong",
		NewNum: newNum,
	}

	serialized_info, err := config.Jsoner.Marshal(hbInfo)
	if err != nil {
		logger.Exceptions.Panic(err)
	}

	return serialized_info
}

func ParseHeartbeatResponse(raw []byte) (protocols.HeartBeatResponse, error) {
	rsp := protocols.HeartBeatResponse{}
	err := config.Jsoner.Unmarshal(raw, &rsp)
	if err != nil {
		logger.Network.Print(err)
		return rsp, err
	}
	return rsp, nil
}
