package config

type ConfigModel struct {
	Name            string            `mapstructure:"name"`
	HttpsNameServer HttpsNsModel      `mapstructure:"httpsNameServer"`
	TentacleFace    TentacleFaceModel `mapstructure:"tentacleFace"`
	OctlFace        OctlFaceModel     `mapstructure:"octlFace"`
	Redis           RedisModel        `mapstructure:"redis"`
	Logger          LoggerModel       `mapstructure:"logger"`
	Workspace       WorkspaceModel    `mapstructure:"workspace"`
	Sslinfo         SslinfoModel      `mapstructure:"ssl"`
	MailAlert       MailAlertModel    `mapstructure:"mailAlert"`
	PakmaServer     PakmaModel        `mapstructure:"pakma"`
	JsonFast        bool              `mapstructure:"jsonFast"`
}

type TentacleFaceModel struct {
	NetDevice     string `mapstructure:"netDevice"`
	Ip            string `mapstructure:"ip"`
	HeartbeatPort uint16 `mapstructure:"heartbeatPort"`
	MessagePort   uint16 `mapstructure:"messagePort"`
	Token         string `mapstructure:"token"`
	ActiveTimeout int    `mapstructure:"activeTimeout"`
	RecordTimeout int    `mapstructure:"recordTimeout"`
}

type OctlFaceModel struct {
	NetDevice string `mapstructure:"netDevice"`
	Ip        string `mapstructure:"ip"`
	Port      uint16 `mapstructure:"port"`
}

type RedisModel struct {
	Ip       string `mapstructure:"ip"`
	Port     uint16 `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Db       int    `mapstructure:"db"`
}

type LoggerModel struct {
	Path       string `mapstructure:"path"`
	NamePrefix string `mapstructure:"namePrefix"`
	RollDays   int    `mapstructure:"rollDays"`
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

type MailAlertModel struct {
	Enabled    bool     `mapstructure:"enable"`
	SmtpServer string   `mapstructure:"smtpServer"`
	SmtpPort   uint16   `mapstructure:"smtpServerPort"`
	Sender     string   `mapstructure:"sender"`
	Password   string   `mapstructure:"password"`
	Receivers  []string `mapstructure:"receivers"`
}

type PakmaModel struct {
	Root            string `mapstructure:"root"`
	Port            uint16 `mapstructure:"port"`
	PreviewDuration int    `mapstructure:"previewDuration"`
}
