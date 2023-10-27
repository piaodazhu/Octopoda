package protocols

type NodeJoinInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Addr    string `json:"addr"`
}

type NodeJoinResponse struct {
	Ts  int64  `json:"ts"`
	NewNum uint32 `json:"new_num"`
}

type HeartBeatRequest struct {
	Msg string `json:"msg"` // reserved for future usage
	Num uint32 `json:"num"`
}

type HeartBeatResponse struct {
	Msg    string `json:"msg"` // reserved for future usage
	NewNum uint32 `json:"new_num"`
}
