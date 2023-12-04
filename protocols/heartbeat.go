package protocols

type NodeJoinInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Addr    string `json:"addr"`
}

type NodeJoinResponse struct {
	Ts     int64  `json:"ts"`
	NewNum uint32 `json:"new_num"`
}

type HeartBeatRequest struct {
	Num   uint32 `json:"num"`
	Delay int64  `json:"delay"`
}

type HeartBeatResponse struct {
	NewNum      uint32 `json:"new_num"`
	Ts          int64  `json:"ts"`
	IsMsgConnOn bool   `json:"is_msgconn_on"`
}
