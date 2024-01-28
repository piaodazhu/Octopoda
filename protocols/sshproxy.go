package protocols

type SSHInfo struct {
	Ip       string
	Port     uint32
	Username string
	Password string
}

type SSHInfoDump struct {
	Name     string
	Username string
	Password string
	Ip       string
	Port     uint32
}

type ProxyMsg struct {
	Code int
	Msg  string
	Data string
}
