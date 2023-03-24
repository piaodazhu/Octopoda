package config


type ConfigModel struct {
	Server ServerModel `mapstructure:"server"`
	Api    APIs        `mapstructure:"api"`
}

type ServerModel struct {
	Name      string `mapstructure:"name"`
	Ip        string `mapstructure:"ip"`
	Port      uint16 `mapstructure:"port"`
	ApiPrefix string `mapstructure:"apiPrefix"`
}

type APIs struct {
	NodeInfo   string `mapstructure:"nodeInfo"`
	NodeStatus string `mapstructure:"nodeStatus"`
	NodeApps   string `mapstructure:"nodeApps"`
	NodeLog    string `mapstructure:"nodeLog"`
	NodeReboot string `mapstructure:"nodeReboot"`

	NodesInfo  string `mapstructure:"nodesInfo"`
	NodesStaus string `mapstructure:"nodesStaus"`

	ScenarioInfo     string `mapstructure:"scenarioInfo"`
	ScenarioVersions string `mapstructure:"scenarioVersions"`
	ScenarioLog      string `mapstructure:"scenarioLog"`

	FileUpload string `mapstructure:"fileUpload"`
	FileSpread string `mapstructure:"fileSpread"`
	FileTree   string `mapstructure:"fileTree"`

	SshInfo string `mapstructure:"sshInfo"`
}
