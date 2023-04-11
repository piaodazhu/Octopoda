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
	NodePrune  string `mapstructure:"nodePrune"`

	NodesInfo   string `mapstructure:"nodesInfo"`
	NodesStatus string `mapstructure:"nodesStatus"`

	ScenarioInfo     string `mapstructure:"scenarioInfo"`
	ScenarioVersion string `mapstructure:"scenarioVersion"`
	ScenarioUpdate string `mapstructure:"scenarioUpdate"`
	ScenarioLog      string `mapstructure:"scenarioLog"`
	ScenarioAppCreate string `mapstructure:"scenarioAppCreate"`
	ScenarioAppDepoly      string `mapstructure:"scenarioAppDeploy"`

	ScenariosInfo     string `mapstructure:"scenariosInfo"`

	FileUpload  string `mapstructure:"fileUpload"`
	FileSpread  string `mapstructure:"fileSpread"`
	FileDistrib string `mapstructure:"fileDistrib"`
	FileTree    string `mapstructure:"fileTree"`

	RunCmd    string `mapstructure:"runCmd"`
	RunScript string `mapstructure:"runScript"`

	SshInfo string `mapstructure:"sshInfo"`
}
