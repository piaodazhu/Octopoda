package config

type ConfigModel struct {
	HttpsNameServer HttpsNsModel `mapstructure:"httpsNameServer"`
	Brain           BrainModel   `mapstructure:"brain"`
	Sslinfo         SslinfoModel `mapstructure:"ssl"`
	Api             APIs         `mapstructure:"api"`
	OutputPretty    bool         `mapstructure:"outputPretty"`
	JsonFast        bool         `mapstructure:"jsonFast"`
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
	Pakma     string `mapstructure:"pakma"`

	TaskState string `mapstructure:"taskState"`

	SshInfo string `mapstructure:"sshInfo"`
}
