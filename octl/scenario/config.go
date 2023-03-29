package scenario

type ScenarioConfigModel struct {
	Name string
	Description string
	Applications []ApplicationConfigModel
}

type ApplicationConfigModel struct {
	Name string
	Description string 
	ScriptPath string
	Nodes []string 
	Prepare ScriptConfigModel
	Run  ScriptConfigModel
	Stop ScriptConfigModel
	Purge ScriptConfigModel
}

type ScriptConfigModel struct {
	Script string
	Files []string
}