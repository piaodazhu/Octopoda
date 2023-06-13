package config

type ConfigModel struct {
	Name string `mapstructure:"name"`
	// NetDevice       string         `mapstructure:"netDevice"`
	// Ip              string         `mapstructure:"-"` // don't need to config
	// Port            uint16         `mapstructure:"port"`
	HttpsNameServer HttpsNsModel `mapstructure:"httpsNameServer"`
	// Worker          WorkerModel    `mapstructure:"tentacle"`
	Sshinfo   SshinfoModel   `mapstructure:"ssh"`
	Brain     BrainModel     `mapstructure:"brain"`
	Logger    LoggerModel    `mapstructure:"logger"`
	Heartbeat HeartbeatModel `mapstructure:"heartbeat"`
	Workspace WorkspaceModel `mapstructure:"workspace"`
	JsonFast  bool           `mapstructure:"jsonFast"`
}

// deprecate
// type WorkerModel struct {
// 	Id   int    `mapstructure:"id"`
// 	Name string `mapstructure:"name"`
// 	Ip   string `mapstructure:"ip"`
// 	Port uint16 `mapstructure:"port"`
// }

type SshinfoModel struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Ip       string `mapstructure:"ip"`
	Port     uint16 `mapstructure:"port"`
}

type BrainModel struct {
	Name          string `mapstructure:"name"`
	Ip            string `mapstructure:"ip"`
	HeartbeatPort uint16 `mapstructure:"heartbeatPort"`
	MessagePort   uint16 `mapstructure:"messagePort"`
	// reserved
	Token string `mapstructure:"token"`
}

type LoggerModel struct {
	Path       string `mapstructure:"path"`
	NamePrefix string `mapstructure:"namePrefix"`
	RollDays   int    `mapstructure:"rollDays"`
}

type HeartbeatModel struct {
	SendInterval      int  `mapstructure:"sendInterval"`
	ReconnectInterval int  `mapstructure:"reconnectInterval"`
	RetryTime         int  `mapstructure:"retryTime"`
	AutoRestart       bool `mapstructure:"autoRestart"`
}

type WorkspaceModel struct {
	Root  string `mapstructure:"root"`
	Store string `mapstructure:"store"`
}

type HttpsNsModel struct {
	Enabled         bool   `mapstructure:"enable"`
	Host            string `mapstructure:"host"`
	Port            uint16 `mapstructure:"port"`
	RequestInterval int    `mapstructure:"requestInterval"`
}
