package heartbeat

import (
	"encoding/hex"
	"github.com/piaodazhu/Octopoda/tentacle/config"
	"github.com/piaodazhu/Octopoda/tentacle/logger"
	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/buildinfo"
)

func MakeNodeJoin(addr string) []byte {
	nodeJoinInfo := protocols.NodeJoinInfo{
		Name:    config.GlobalConfig.Name,
		Version: buildinfo.String(),
		Addr:    addr,
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

func ParseNodeJoinResponse(raw []byte) (protocols.NodeJoinResponse, error) {
	decBuffer := make([]byte, hex.DecodedLen(len(raw)))
	_, err := hex.Decode(decBuffer, raw)
	if err != nil {
		logger.Exceptions.Print(err)
		return protocols.NodeJoinResponse{}, err
	}

	buffer := decBuffer
	response := protocols.NodeJoinResponse{}
	err = config.Jsoner.Unmarshal(buffer, &response)
	if err != nil {
		logger.Exceptions.Print(err)
		return protocols.NodeJoinResponse{}, err
	}
	return response, nil
}
