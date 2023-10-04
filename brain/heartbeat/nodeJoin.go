package heartbeat

import (
	"brain/config"
	"brain/logger"
	"encoding/hex"
	"time"
)

type NodeJoinInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Addr    string `json:"addr"`
}

type NodeJoinResponse struct {
	Ts int64 `json:"ts"`
}

func ParseNodeJoin(raw []byte) (NodeJoinInfo, error) {
	decBuffer := make([]byte, hex.DecodedLen(len(raw)))
	_, err := hex.Decode(decBuffer, raw)
	if err != nil {
		logger.Network.Print(err)
		return NodeJoinInfo{}, err
	}

	buffer := decBuffer

	info := NodeJoinInfo{}
	err = config.Jsoner.Unmarshal(buffer, &info)
	if err != nil {
		logger.Network.Print(err)
		return NodeJoinInfo{}, err
	}
	return info, nil
}

func MakeNodeJoinResponse() []byte {
	nodeJoinResponse := NodeJoinResponse{
		Ts: time.Now().UnixMicro(),
	}
	serialized_response, err := config.Jsoner.Marshal(nodeJoinResponse)
	if err != nil {
		logger.Exceptions.Panic(err)
	}

	buffer := serialized_response

	encBuffer := make([]byte, hex.EncodedLen(len(buffer)))
	hex.Encode(encBuffer, buffer)
	return encBuffer
}
