package heartbeat

import (
	"time"

	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/ostp"
	"github.com/piaodazhu/Octopoda/tentacle/config"
)

func MakeHeartbeat(num uint32) []byte {
	hbInfo := protocols.HeartBeatRequest{
		Num:   num,
		Delay: ostp.Delay,
	}

	serialized_info, _ := config.Jsoner.Marshal(hbInfo)
	return serialized_info
}

func ParseHeartbeat(raw []byte) (protocols.HeartBeatRequest, error) {
	info := protocols.HeartBeatRequest{}
	err := config.Jsoner.Unmarshal(raw, &info)
	if err != nil {
		return info, err
	}
	return info, nil
}

func MakeHeartbeatResponse(newNum uint32, isMsgConnOn bool) []byte {
	hbInfo := protocols.HeartBeatResponse{
		NewNum:      newNum,
		Ts:          time.Now().UnixMilli(),
		IsMsgConnOn: isMsgConnOn,
	}

	serialized_info, _ := config.Jsoner.Marshal(hbInfo)
	return serialized_info
}

func ParseHeartbeatResponse(raw []byte) (protocols.HeartBeatResponse, error) {
	rsp := protocols.HeartBeatResponse{}
	err := config.Jsoner.Unmarshal(raw, &rsp)
	if err != nil {
		return rsp, err
	}
	return rsp, nil
}
