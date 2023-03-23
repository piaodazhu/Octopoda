package model

const (
	NodeStateReady = iota
	NodeStateDisconn
	NodeStateDead
)

type NodeModel struct {
	Id        uint32
	Name      string
	Addr      string
	State     int32
	OnlineTs  int64
	OfflineTs int64
	Applist   []*AppModel
}

type AppModel struct {
	Id          uint32
	Name        string
	Description string
	Versions    []*AppVersionModel
	Scenario    *ScenarioModel
	Node        *NodeModel
}

type AppVersionModel struct {
	Version   uint64
	Message   string
	Timestamp int64
	// Tag     string
}

type ScenarioVersionModel struct {
	AppVersionModel
	Apps []*AppModel
}

type ScenarioModel struct {
	Id          uint32
	Name        string
	Description string
	Versions    []*ScenarioVersionModel
}

type State struct {
	Id        int
	Name      string
	Platform  string
	CpuCores  int
	Ip        string
	StartTime int64

	CpuLoadShort float64
	CpuLoadLong  float64
	MemUsed      uint64
	MemTotal     uint64
	DiskUsed     uint64
	DiskTotal    uint64
}