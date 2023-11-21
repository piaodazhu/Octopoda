package config

type ConfigModel struct {
	Name            string           `mapstructure:"name"`
	NetDevice       string           `mapstructure:"netDevice"`
	HttpsNameServer HttpsNsModel     `mapstructure:"httpsNameServer"`
	Brain           BrainModel       `mapstructure:"brain"`
	Logger          LoggerModel      `mapstructure:"logger"`
	Heartbeat       HeartbeatModel   `mapstructure:"heartbeat"`
	Workspace       WorkspaceModel   `mapstructure:"workspace"`
	Sslinfo         SslinfoModel     `mapstructure:"ssl"`
	CustomEnv       []*CustomEnvItem `mapstructure:"env"`
	PakmaServer     PakmaModel       `mapstructure:"pakma"`
	JsonFast        bool             `mapstructure:"jsonFast"`
}

type BrainModel struct {
	Name          string `mapstructure:"name"`
	Ip            string `mapstructure:"ip"`
	HeartbeatPort uint16 `mapstructure:"heartbeatPort"`
	MessagePort   uint16 `mapstructure:"messagePort"`
}

type LoggerModel struct {
	Path       string `mapstructure:"path"`
	NamePrefix string `mapstructure:"namePrefix"`
	RollDays   int    `mapstructure:"rollDays"`
}

type HeartbeatModel struct {
	SendInterval       int    `mapstructure:"sendInterval"`
	ReconnectInterval  int    `mapstructure:"reconnectInterval"`
	RetryTime          int    `mapstructure:"retryTime"`
	AutoRestart        bool   `mapstructure:"autoRestart"`
	AutoRestartCommand string `mapstructure:"autoRestartCommand"`
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

type SslinfoModel struct {
	CaCert     string `mapstructure:"caCert"`
	ClientCert string `mapstructure:"clientCert"`
	ClientKey  string `mapstructure:"clientKey"`
}

type CustomEnvItem struct {
	Key   string `mapstructure:"key"`
	Value string `mapstructure:"value"`
}

type PakmaModel struct {
	Root            string `mapstructure:"root"`
	Port            uint16 `mapstructure:"port"`
	PreviewDuration int    `mapstructure:"previewDuration"`
}
