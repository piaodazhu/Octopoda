package san

// 1 model
// ScenarioVersion: a snapshot of a scenario
type ScenarioVersionModel struct {
	Version
	Apps []*AppModel
}

// App: a kind of application in a scenario
type AppModel struct {
	Id          uint32
	Name        string
	Description string
	NodeApp     []*NodeAppModel
}

// NodeApp: a application instence on the node
type NodeAppModel struct {
	Name    string
	Version string
}

// 2 params
// digest of a scenario
type ScenarioDigest struct {
	Name        string
	Description string
	Version     string
	Timestamp   TsInt64
	Message     string
}

// detail info of a scenario
type ScenarioInfo struct {
	ScenarioDigest
	Apps []*AppInfo
}
type AppInfo struct {
	Name        string
	Description string
	NodeApps    []string
}

type NodeAppItem struct {
	AppName  string
	ScenName string
	NodeName string
	Version  string
}
