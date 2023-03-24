package config

type ConfigModel struct {
	TentacleFace TentacleFaceModel `mapstructure:"tentacleFace"`
	BrainFace    BrainFaceModel    `mapstructure:"brainFace"`
	Sshinfo      SshinfoModel      `mapstructure:"ssh"`
	Redis        RedisModel        `mapstructure:"redis"`
	Logger       LoggerModel       `mapstructure:"logger"`
	Workspace    WorkspaceModel    `mapstructure:"workspace"`
}

type TentacleFaceModel struct {
	Ip            string `mapstructure:"ip"`
	Port          uint16 `mapstructure:"port"`
	Token         string `mapstructure:"token"`
	ActiveTimeout int    `mapstructure:"activeTimeout"`
	RecordTimeout int    `mapstructure:"recordTimeout"`
}

type BrainFaceModel struct {
	Ip   string `mapstructure:"ip"`
	Port uint16 `mapstructure:"port"`
}

type SshinfoModel struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Ip       string `mapstructure:"ip"`
	Port     uint16 `mapstructure:"port"`
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
	Root string `mapstructure:"root"`
}
