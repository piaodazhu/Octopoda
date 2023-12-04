package config

type ConfigModel struct {
	HttpsNameServer HttpsNsModel   `mapstructure:"httpsNameServer"`
	Brain           BrainModel     `mapstructure:"brain"`
	Sslinfo         SslinfoModel   `mapstructure:"ssl"`
	Gitinfo         GitinfoModel   `mapstructure:"git"`
	OutputPretty    bool           `mapstructure:"outputPretty"`
	JsonFast        bool           `mapstructure:"jsonFast"`
	Workgroup       WorkgroupModel `mapstructure:"workgroup"`
}

type BrainModel struct {
	Name      string `mapstructure:"name"`
	Ip        string `mapstructure:"ip"`
	Port      uint16 `mapstructure:"port"`
	ApiPrefix string `mapstructure:"apiPrefix"`
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

type GitinfoModel struct {
	ServeUrl string `mapstructure:"serveUrl"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type WorkgroupModel struct {
	Root            string `mapstructure:"root"`
	Password        string `mapstructure:"password"`
	CurrentPathFile string `mapstructure:"currentPathFile"`
}
