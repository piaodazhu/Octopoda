package heartbeat

import (
	"brain/config"
	"brain/logger"
	"brain/ticker"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
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
	key := tokenToKey(config.GlobalConfig.TentacleFace.Token, 16)
	blk, err := aes.NewCipher(key)
	if err != nil {
		logger.Tentacle.Panic(err)
	}
	block = blk
}

func ParseNodeJoin(raw []byte) (NodeJoinInfo, error) {
	decBuffer := make([]byte, hex.DecodedLen(len(raw)))
	_, err := hex.Decode(decBuffer, raw)
	if err != nil {
		logger.Tentacle.Print(err)
		return NodeJoinInfo{}, err
	}
	// logger.Tentacle.Print(raw, decBuffer)

	// buffer := make([]byte, len(decBuffer))
	// block.Decrypt(buffer, decBuffer)
	buffer := decBuffer

	info := NodeJoinInfo{}
	err = json.Unmarshal(buffer, &info)
	if err != nil {
		logger.Tentacle.Print(err)
		return NodeJoinInfo{}, err
	}
	return info, nil
}

func MakeNodeJoinResponse() []byte {
	nodeJoinResponse := NodeJoinResponse{
		Ts:  time.Now().Unix(),
		Cnt: ticker.GetTick(),
	}
	serialized_response, err := json.Marshal(nodeJoinResponse)
	if err != nil {
		logger.Tentacle.Panic(err)
	}

	// buffer := make([]byte, len(serialized_response))
	// block.Encrypt(buffer, serialized_response)
	buffer := serialized_response

	encBuffer := make([]byte, hex.EncodedLen(len(buffer)))
	hex.Encode(encBuffer, buffer)
	return encBuffer
}
