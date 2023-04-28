package config

type ConfigModel struct {
	Server ServerModel `mapstructure:"server"`
	Api    APIs        `mapstructure:"api"`
	OutputPretty bool   `mapstructure:"outputPretty"`
	JsonFast     bool   `mapstructure:"jsonFast"`
}

type ServerModel struct {
	Name         string `mapstructure:"name"`
	Ip           string `mapstructure:"ip"`
	Port         uint16 `mapstructure:"port"`
	ApiPrefix    string `mapstructure:"apiPrefix"`
}

type APIs struct {
	NodeInfo       string `mapstructure:"nodeInfo"`
	NodeStatus     string `mapstructure:"nodeStatus"`
	NodeLog        string `mapstructure:"nodeLog"`
	NodeReboot     string `mapstructure:"nodeReboot"`
	NodePrune      string `mapstructure:"nodePrune"`
	NodeApps       string `mapstructure:"nodeApps"`
	NodeAppVersion string `mapstructure:"nodeAppVersion"`
	NodeAppReset   string `mapstructure:"nodeAppReset"`

	NodesInfo   string `mapstructure:"nodesInfo"`
	NodesStatus string `mapstructure:"nodesStatus"`

	ScenarioInfo      string `mapstructure:"scenarioInfo"`
	ScenarioVersion   string `mapstructure:"scenarioVersion"`
	ScenarioUpdate    string `mapstructure:"scenarioUpdate"`
	ScenarioFix       string `mapstructure:"scenarioFix"`
	ScenarioAppCreate string `mapstructure:"scenarioAppCreate"`
	ScenarioAppDepoly string `mapstructure:"scenarioAppDeploy"`

	ScenariosInfo string `mapstructure:"scenariosInfo"`

	FileUpload  string `mapstructure:"fileUpload"`
	FileSpread  string `mapstructure:"fileSpread"`
	FileDistrib string `mapstructure:"fileDistrib"`
	FileTree    string `mapstructure:"fileTree"`
	FilePull    string `mapstructure:"filePull"`

	RunCmd    string `mapstructure:"runCmd"`
	RunScript string `mapstructure:"runScript"`

	TaskState string `mapstructure:"taskState"`

	SshInfo string `mapstructure:"sshInfo"`
}
