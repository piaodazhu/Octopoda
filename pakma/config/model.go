package config

type ConfigModel struct {
	Packma          PackmaModel  `mapstructure:"pakma"`
	HttpsNameServer HttpsNsModel `mapstructure:"httpsNameServer"`
	Sslinfo         SslinfoModel `mapstructure:"ssl"`
	AppName         string
	AppOS           string
	AppArch         string
}

type PackmaModel struct {
	ServePort       uint16 `mapstructure:"port"`
	Root            string `mapstructure:"root"`
	PreviewDuration int    `mapstructure:"previewDuration"`
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
