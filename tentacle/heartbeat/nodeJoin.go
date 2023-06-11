package heartbeat

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"tentacle/config"
	"tentacle/logger"
	"time"
)

type NodeJoinMessage struct {
	Type int    `json:"type"`
	Raw  string `json:"raw"`
}

type NodeJoinInfo struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	Port uint16 `json:"port"`
	Ts   int64  `json:"ts"`
}

type NodeJoinResponse struct {
	Ts  int64 `json:"ts"`
	Cnt int64 `json:"cnt"`
}

var block cipher.Block

func tokenToKey(key string, targetsize int) []byte {
	keybuf := []byte(key)
	if len(keybuf) >= targetsize {
		return keybuf[:targetsize]
	}
	for i := len(keybuf); i < targetsize; i++ {
		keybuf = append(keybuf, '=')
	}
	return keybuf
}

func InitHeartbeat() {
	key := tokenToKey(config.GlobalConfig.Master.Token, 16)
	blk, err := aes.NewCipher(key)
	if err != nil {
		logger.Exceptions.Panic(err)
	}
	block = blk
}

func MakeNodeJoin() []byte {
	nodeJoinInfo := NodeJoinInfo{
		Name: config.GlobalConfig.Name,
		IP:   config.GlobalConfig.Ip,
		Port: config.GlobalConfig.Port,
		Ts:   time.Now().Unix(),
	}
	serialized_info, err := config.Jsoner.Marshal(nodeJoinInfo)
	if err != nil {
		logger.Exceptions.Panic(err)
		return nil
	}

	// buffer := make([]byte, len(serialized_info))
	// block.Encrypt(buffer, serialized_info)

	// logger.Client.Print(serialized_info, buffer)
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

	// buffer := make([]byte, len(decBuffer))
	// block.Decrypt(buffer, decBuffer)
	buffer := decBuffer

	response := NodeJoinResponse{}
	err = config.Jsoner.Unmarshal(buffer, &response)
	if err != nil {
		logger.Exceptions.Print(err)
		return NodeJoinResponse{}, err
	}
	return response, nil
}
