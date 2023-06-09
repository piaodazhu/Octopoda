package heartbeat

import (
	"encoding/hex"
	"tentacle/config"
	"tentacle/logger"
)

type NodeJoinInfo struct {
	Name string `json:"name"`
}

type NodeJoinResponse struct {
	Ts  int64 `json:"ts"`
}

func MakeNodeJoin() []byte {
	nodeJoinInfo := NodeJoinInfo{
		Name: config.GlobalConfig.Name,
	}
	serialized_info, err := config.Jsoner.Marshal(nodeJoinInfo)
	if err != nil {
		logger.Exceptions.Panic(err)
		return nil
	}

	buffer := serialized_info
	encBuffer := make([]byte, hex.EncodedLen(len(buffer)))
	hex.Encode(encBuffer, buffer)
	return encBuffer
}

func ParseNodeJoinResponse(raw []byte) (NodeJoinResponse, error) {
	decBuffer := make([]byte, hex.DecodedLen(len(raw)))
	_, err := hex.Decode(decBuffer, raw)
	if err != nil {
		logger.Exceptions.Print(err)
		return NodeJoinResponse{}, err
	}

	buffer := decBuffer
	response := NodeJoinResponse{}
	err = config.Jsoner.Unmarshal(buffer, &response)
	if err != nil {
		logger.Exceptions.Print(err)
		return NodeJoinResponse{}, err
	}
	return response, nil
}
