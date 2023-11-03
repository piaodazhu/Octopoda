package protocols

import (
	"fmt"
	"sort"
	"time"
)

const (
	NodeStateReady = iota
	NodeStateDisconn
	NodeStateDead
)

type NodeInfo struct {
	Name      string
	Version   string
	Addr      string
	State     int32
	ConnState string
	OnlineTs  int64
	OfflineTs int64
	ActiveTs  int64
}

type BrainInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Addr    string `json:"addr"`
}

type NodesInfo struct {
	BrainName    string
	BrainVersion string
	BrainAddr    string
	InfoList     []NodeInfo
}

type Status struct {
	Name      string
	Platform  string
	CpuCores  int
	LocalTime time.Time

	CpuLoadShort float64
	CpuLoadLong  float64
	MemUsed      uint64
	MemTotal     uint64
	DiskUsed     uint64
	DiskTotal    uint64
}

type NodesStatus struct {
	StatusList []Status
}

type NodeInfoText struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Addr         string `json:"addr"`
	Health       string `json:"health"`
	MsgConnState string `json:"msg_conn"`
	OnlineTime   string `json:"online_time,omitempty"`
	OfflineTime  string `json:"offline_time,omitempty"`
	LastOnline   string `json:"last_active,omitempty"`
}

type NodesInfoText struct {
	BrainInfo    BrainInfoText   `json:"brain"`
	NodeInfoList []*NodeInfoText `json:"nodes"`
	Total        int             `json:"total"`
	Active       int             `json:"active"`
	Offline      int             `json:"offline"`
}

type BrainInfoText struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Addr    string `json:"addr"`
}

type StatusText struct {
	Name         string `json:"name"`
	Platform     string `json:"platform"`
	CpuCores     int    `json:"cpu_cores"`
	LocalTime    string `json:"local_time"`
	CpuLoadShort string `json:"cpu_average1"`
	CpuLoadLong  string `json:"cpu_average10"`
	MemUsage     string `json:"memory_usage"`
	DiskUsage    string `json:"disk_usage"`
}

type NodesStatusText struct {
	NodesStatusList []*StatusText `json:"nodes"`
	AvrCpuLoad      string        `json:"average_cpuload"`
	AvrMemoryUsage  string        `json:"average_memoryusage"`
	AvrDiskUsage    string        `json:"average_diskusage"`
}

func (node *NodeInfo) ToText() *NodeInfoText {
	res := &NodeInfoText{
		Name:    node.Name,
		Version: node.Version,
		Addr:    node.Addr,
	}
	switch node.State {
	case 0:
		res.Health = "Healthy"
		res.OnlineTime = time.Since(time.UnixMilli(node.OnlineTs)).String()
	case 1:
		res.Health = "Disconnect"
		res.LastOnline = time.UnixMilli(node.ActiveTs).Format("2006-01-02 15:04:05")
	case 2:
		res.Health = "Offline"
		res.OfflineTime = time.Since(time.UnixMilli(node.OfflineTs)).String()
	}
	res.MsgConnState = node.ConnState
	return res
}

func (nodes *NodesInfo) ToText() *NodesInfoText {
	res := &NodesInfoText{
		BrainInfo: BrainInfoText{
			Name:    nodes.BrainName,
			Version: nodes.BrainVersion,
			Addr:    nodes.BrainAddr,
		},
		NodeInfoList: make([]*NodeInfoText, len(nodes.InfoList)),
	}
	for i, node := range nodes.InfoList {
		res.Total++
		if node.State == 0 {
			res.Active++
		} else if node.State == 2 || node.ConnState != "On" {
			res.Offline++
		}
		res.NodeInfoList[i] = node.ToText()
	}
	sort.Slice(res.NodeInfoList, func(i, j int) bool {
		return res.NodeInfoList[i].Name < res.NodeInfoList[j].Name
	})
	return res
}

func (status *Status) ToText() *StatusText {
	return &StatusText{
		Name:         status.Name,
		Platform:     status.Platform,
		CpuCores:     status.CpuCores,
		LocalTime:    status.LocalTime.Format("2006-01-02 15:04:05"),
		CpuLoadShort: fmt.Sprintf("%5.1f%%", status.CpuLoadShort),
		CpuLoadLong:  fmt.Sprintf("%5.1f%%", status.CpuLoadLong),
		MemUsage: fmt.Sprintf("%5.1f%%: (%.2fGB / %.2fGB)",
			float64(status.MemUsed*100)/float64(status.MemTotal),
			float64(status.MemUsed)/1073741824,
			float64(status.MemTotal)/1073741824),
		DiskUsage: fmt.Sprintf("%5.1f%%: (%.2fGB / %.2fGB)",
			float64(status.DiskUsed*100)/float64(status.DiskTotal),
			float64(status.DiskUsed)/1073741824,
			float64(status.DiskTotal)/1073741824),
	}
}

func (status *NodesStatus) ToText() *NodesStatusText {
	var nodesStatus NodesStatusText
	var cpu_load_sum float64 = 0.0
	var mem_used_sum, mem_tot_sum uint64 = 0, 0
	var disk_used_sum, disk_tot_sum uint64 = 0, 0

	for _, status := range status.StatusList {
		nodesStatus.NodesStatusList = append(nodesStatus.NodesStatusList, status.ToText())
		cpu_load_sum += status.CpuLoadLong
		mem_used_sum += status.MemUsed
		mem_tot_sum += status.MemTotal
		disk_used_sum += status.DiskUsed
		disk_tot_sum += status.DiskTotal
	}
	nodesStatus.AvrCpuLoad = fmt.Sprintf("%5.1f%%", cpu_load_sum/float64(len(status.StatusList)))
	nodesStatus.AvrMemoryUsage = fmt.Sprintf("%5.1f%%", float64(mem_used_sum*100)/float64(mem_tot_sum))
	nodesStatus.AvrDiskUsage = fmt.Sprintf("%5.1f%%", float64(disk_used_sum*100)/float64(disk_tot_sum))

	sort.Slice(nodesStatus.NodesStatusList, func(i, j int) bool {
		return nodesStatus.NodesStatusList[i].Name < nodesStatus.NodesStatusList[j].Name
	})
	return &nodesStatus
}
