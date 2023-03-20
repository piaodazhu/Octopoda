package config

type ConfigModel struct {
	TentacleFace TentacleFaceModel `mapstructure:"tentacleFace"`
	BrainFace    BrainFaceModel    `mapstructure:"brainFace"`
	Redis        RedisModel        `mapstructure:"redis"`
	Logger       LoggerModel       `mapstructure:"logger"`
	Workspace    WorkspaceModel    `mapstructure:"workspace"`
}

type TentacleFaceModel struct {
	Ip            int    `mapstructure:"ip"`
	Port          uint16 `mapstructure:"port"`
	Token         string `mapstructure:"token"`
	CheckInterval int    `mapstructure:"checkInterval"`
	PurgeMarkTime int    `mapstructure:"purgeMarkTime"`
}

type BrainFaceModel struct {
	Ip   string `mapstructure:"ip"`
	Port uint16 `mapstructure:"port"`
}

type RedisModel struct {
	Ip   string `mapstructure:"ip"`
	Port uint16 `mapstructure:"port"`
	Db   int    `mapstructure:"db"`
}

type LoggerModel struct {
	Path       string `mapstructure:"path"`
	NamePrefix string `mapstructure:"namePrefix"`
	RollDays   int    `mapstructure:"rollDays"`
}

type WorkspaceModel struct {
	Root string `mapstructure:"root"`
}
