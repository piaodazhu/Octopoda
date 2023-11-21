package heartbeat

import (
	"time"

	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/logger"
	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/ostp"
)

func MakeHeartbeat(num uint32) []byte {
	hbInfo := protocols.HeartBeatRequest{
		Msg:   "ping",
		Num:   num,
		Delay: ostp.Delay,
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
		Ts:     time.Now().UnixMilli(),
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
