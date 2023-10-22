package config

const (
	API_NodeInfo       = "/node/info"
	API_NodeStatus     = "/node/status"
	API_NodeLog        = "/node/log"
	API_NodePrune      = "/node/prune"
	API_NodeApps       = "/node/apps"
	API_NodeAppVersion = "/node/app/version"
	API_NodeAppReset   = "/node/app/version"

	API_NodesInfo   = "/nodes/info"
	API_NodesStatus = "/nodes/status"
	API_NodesParse  = "/nodes/parse"

	API_ScenarioInfo      = "/scenario/info"
	API_ScenarioUpdate    = "/scenario/update"
	API_ScenarioVersion   = "/scenario/version"
	API_ScenarioFix       = "/scenario/fix"
	API_ScenarioAppDeploy = "/scenario/app/deployment"
	API_ScenarioAppCreate = "/scenario/app/prepare"

	API_ScenariosInfo = "/scenarios/info"

	API_FileUpload  = "/file/upload"
	API_FileSpread  = "/file/spread"
	API_FileDistrib = "/file/distrib"
	API_FileTree    = "/file/tree"
	API_FilePull    = "/file/pull"

	API_RunCmd    = "/run/cmd"
	API_RunScript = "/run/script"
	API_RunCancel = "/run/cancel"

	API_Ssh = "/ssh"

	API_TaskState = "/taskstate"

	API_Pakma = "/pakma"

	API_Group  = "/group"
	API_Groups = "/groups"
)
