package config

type ConfigModel struct {
	Worker    WorkerModel    `mapstructure:"tentacle"`
	Sshinfo   SshinfoModel   `mapstructure:"ssh"`
	Master    MasterModel    `mapstructure:"brain"`
	Logger    LoggerModel    `mapstructure:"logger"`
	Heartbeat HeartbeatModel `mapstructure:"heartbeat"`
	Workspace WorkspaceModel `mapstructure:"workspace"`
}

type WorkerModel struct {
	Id   int    `mapstructure:"id"`
	Name string `mapstructure:"name"`
	Ip   string `mapstructure:"ip"`
	Port uint16 `mapstructure:"port"`
}

type SshinfoModel struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Ip       string `mapstructure:"ip"`
	Port     uint16 `mapstructure:"port"`
}

type MasterModel struct {
	Ip    string `mapstructure:"ip"`
	Port  uint16 `mapstructure:"port"`
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
	Root string `mapstructure:"root"`
}
