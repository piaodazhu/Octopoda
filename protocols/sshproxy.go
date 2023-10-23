package protocols


type SSHInfo struct {
	Ip       string
	Port     uint32
	Username string
	Password string
}

type ProxyMsg struct {
	Code int
	Msg  string
	Data string
}
