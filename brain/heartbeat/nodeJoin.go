package heartbeat

import (
	"encoding/hex"
	"time"

	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/brain/logger"
	"github.com/piaodazhu/Octopoda/protocols"
)

func ParseNodeJoin(raw []byte) (protocols.NodeJoinInfo, error) {
	decBuffer := make([]byte, hex.DecodedLen(len(raw)))
	_, err := hex.Decode(decBuffer, raw)
	if err != nil {
		logger.Network.Print(err)
		return protocols.NodeJoinInfo{}, err
	}

	buffer := decBuffer

	info := protocols.NodeJoinInfo{}
	err = config.Jsoner.Unmarshal(buffer, &info)
	if err != nil {
		logger.Network.Print(err)
		return protocols.NodeJoinInfo{}, err
	}
	return info, nil
}

func MakeNodeJoinResponse(randNum uint32) []byte {
	nodeJoinResponse := protocols.NodeJoinResponse{
		Ts:     time.Now().UnixMicro(),
		NewNum: randNum,
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
