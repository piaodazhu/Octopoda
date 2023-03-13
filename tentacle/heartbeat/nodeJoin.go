package heartbeat

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"nworkerd/config"
	"nworkerd/logger"
	"nworkerd/message"
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
		logger.Client.Panic(err)
	}
	block = blk
}

func MakeNodeJoin() []byte {
	nodeJoinInfo := NodeJoinInfo{
		Name: config.GlobalConfig.Worker.Name,
		IP:   config.GlobalConfig.Worker.Ip,
		Port: config.GlobalConfig.Worker.Port,
		Ts:   time.Now().Unix(),
	}
	serialized_info, err := json.Marshal(nodeJoinInfo)
	if err != nil {
		logger.Client.Panic(err)
		return nil
	}

	buffer := make([]byte, len(serialized_info))
	block.Encrypt(buffer, serialized_info)

	encBuffer := make([]byte, hex.EncodedLen(len(buffer)))
	hex.Encode(encBuffer, buffer)
	return encBuffer
}

func ParseNodeJoin(raw []byte) (NodeJoinInfo, error) {
	decBuffer := make([]byte, hex.DecodedLen(len(raw)))
	_, err := hex.Decode(decBuffer, raw)
	if err != nil {
		logger.Client.Print(err)
		return NodeJoinInfo{}, err
	}

	buffer := make([]byte, len(decBuffer))
	block.Decrypt(buffer, decBuffer)

	info := NodeJoinInfo{}
	err = json.Unmarshal(buffer, &info)
	if err != nil {
		logger.Client.Print(err)
		return NodeJoinInfo{}, err
	}
	return info, nil
}

func MakeNodeJoinResponse() []byte {
	nodeJoinResponse := NodeJoinResponse{
		Ts:  time.Now().Unix(),
		Cnt: message.GetVersion(),
	}
	serialized_response, err := json.Marshal(nodeJoinResponse)
	if err != nil {
		logger.Client.Panic(err)
	}

	buffer := make([]byte, len(serialized_response))
	block.Encrypt(buffer, serialized_response)

	encBuffer := make([]byte, hex.EncodedLen(len(buffer)))
	hex.Encode(encBuffer, buffer)
	return encBuffer
}

func ParseNodeJoinResponse(raw []byte) (NodeJoinResponse, error) {
	decBuffer := make([]byte, hex.DecodedLen(len(raw)))
	_, err := hex.Decode(decBuffer, raw)
	if err != nil {
		logger.Client.Print(err)
		return NodeJoinResponse{}, err
	}

	buffer := make([]byte, len(decBuffer))
	block.Decrypt(buffer, decBuffer)

	response := NodeJoinResponse{}
	err = json.Unmarshal(buffer, &response)
	if err != nil {
		logger.Client.Print(err)
		return NodeJoinResponse{}, err
	}
	return response, nil
}
