package san

// 1 models
type App struct {
	Name        string
	Discription string
	Versions    []Version
}

type NodeApps struct {
	NodeVersion int64
	Apps        []App
}

type NodeAppsDigest struct {
	Apps []AppDigest
}

type AppDigest struct {
	Name        string
	Discription string
	CurVersion  Version
}

// 2 params
type AppBasic struct {
	Name        string
	Scenario    string
	Description string
	Message     string
}

type AppVersionParams struct {
	AppBasic
	Offset int
	Limit  int
}

type AppCreateParams struct {
	AppBasic
	FilePack string
}

type AppDeployParams struct {
	AppBasic
	Script string
}

type AppDeleteParams struct {
	Name     string
	Scenario string
}

type AppResetParams struct {
	AppBasic
	VersionHash string
	Mode        string
}
