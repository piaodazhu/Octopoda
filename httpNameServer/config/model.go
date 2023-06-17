package config

type ConfigModel struct {
	NetDevice  string `mapstructure:"NetDevice"`
	ServeIp    string
	ServePort  uint16            `mapstructure:"servePort"`
	Redis      RedisModel        `mapstructure:"redis"`
	Logger     LoggerModel       `mapstructure:"logger"`
	Sslinfo    SslinfoModel      `mapstructure:"ssl"`
	StaticDirs []*StaticDirModel `mastructure:"staticDir"`
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

type SslinfoModel struct {
	CaCert     string `mapstructure:"caCert"`
	ServerCert string `mapstructure:"serverCert"`
	ServerKey  string `mapstructure:"serverKey"`
}

type StaticDirModel struct {
	Url string `mapstructure:"url"`
	Dir string `mapstructure:"dir"`
}
