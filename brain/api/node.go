package api

import (
	"brain/config"
	"brain/logger"
	"brain/message"
	"brain/model"
	"brain/sys"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type NodeInfoText struct {
	Name         string `json:"name"`
	Health       string `json:"health"`
	MsgConnState string `json:"msg_conn"`
	OnlineTime   string `json:"online_time,omitempty"`
	OfflineTime  string `json:"offline_time,omitempty"`
	LastOnline   string `json:"last_active,omitempty"`
}

type NodesInfoText struct {
	NodeInfoList []*NodeInfoText `json:"nodes"`
	Total        int             `json:"total"`
	Active       int             `json:"active"`
	Offline      int             `json:"offline"`
}

func nodeInfoToText(node *model.NodeModel) *NodeInfoText {
	res := &NodeInfoText{
		Name: node.Name,
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
	if node.MsgConn == nil {
		res.MsgConnState = "Off"
	} else {
		res.MsgConnState = "On"
	}
	return res
}

func nodesInfoToText(nodes []*model.NodeModel) *NodesInfoText {
	res := &NodesInfoText{
		NodeInfoList: make([]*NodeInfoText, len(nodes)),
	}
	for i, node := range nodes {
		res.Total++
		if node.State == 0 {
			res.Active++
		} else if node.State == 2 {
			res.Offline++
		}
		res.NodeInfoList[i] = nodeInfoToText(node)
	}
	sort.Slice(res.NodeInfoList, func(i, j int) bool {
		return res.NodeInfoList[i].Name < res.NodeInfoList[j].Name
	})
	return res
}

func NodeInfo(ctx *gin.Context) {
	var name string
	var ok bool
	var node *model.NodeModel
	if name, ok = ctx.GetQuery("name"); !ok {
		ctx.JSON(404, struct{}{})
		return
	}
	if node, ok = model.GetNodeInfoByName(name); !ok {
		ctx.JSON(404, struct{}{})
		return
	}
	ctx.JSON(200, nodeInfoToText(node))
}

func NodesInfo(ctx *gin.Context) {
	var nodes []*model.NodeModel
	var ok bool

	if nodes, ok = model.GetNodesInfoAll(); !ok {
		ctx.JSON(404, struct{}{})
		return
	}
	ctx.JSON(200, nodesInfoToText(nodes))
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

func statusToText(status model.Status) StatusText {
	return StatusText{
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

func NodeStatus(ctx *gin.Context) {
	var name string
	var ok bool
	if name, ok = ctx.GetQuery("name"); !ok {
		ctx.JSON(404, struct{}{})
		return
	}
	var status model.Status
	if name == "master" {
		ctx.JSON(200, statusToText(sys.LocalStatus()))
		return
	}
	raw, err := model.Request(name, message.TypeNodeStatus, []byte{})
	if err != nil {
		logger.Comm.Println("NodeStatus", err)
		ctx.JSON(404, struct{}{})
		return
	}
	
	err = json.Unmarshal(raw, &status)
	if err != nil {
		logger.Comm.Println("NodeStatus Unmarshal", err)
		ctx.JSON(404, struct{}{})
		return
	}
	ctx.JSON(200, statusToText(status))
}

func NodesState(ctx *gin.Context) {
	var nodes []*model.NodeModel
	var nodesStatus NodesStatusText
	var ok bool

	if nodes, ok = model.GetNodesInfoAll(); !ok {
		ctx.JSON(404, struct{}{})
	}

	channel := make(chan model.Status, len(nodes))
	var wg sync.WaitGroup
	wg.Add(len(nodes))
	for _, node := range nodes {
		go getNodeStatus(node.Name, channel, &wg)
	}
	wg.Wait()
	close(channel)

	var cpu_load_sum float64 = 0.0
	var mem_used_sum, mem_tot_sum uint64 = 0, 0
	var disk_used_sum, disk_tot_sum uint64 = 0, 0

	for v := range channel {
		text := statusToText(v)
		nodesStatus.NodesStatusList = append(nodesStatus.NodesStatusList, &text)
		cpu_load_sum += v.CpuLoadLong
		mem_used_sum += v.MemUsed
		mem_tot_sum += v.MemTotal
		disk_used_sum += v.DiskUsed
		disk_tot_sum += v.DiskTotal
	}
	nodesStatus.AvrCpuLoad = fmt.Sprintf("%5.1f%%", cpu_load_sum/float64(len(nodes)))
	nodesStatus.AvrMemoryUsage = fmt.Sprintf("%5.1f%%", float64(mem_used_sum*100)/float64(mem_tot_sum))
	nodesStatus.AvrDiskUsage = fmt.Sprintf("%5.1f%%", float64(disk_used_sum*100)/float64(disk_tot_sum))

	sort.Slice(nodesStatus.NodesStatusList, func(i, j int) bool {
		return nodesStatus.NodesStatusList[i].Name < nodesStatus.NodesStatusList[j].Name
	})

	ctx.JSON(200, nodesStatus)
}

func getNodeStatus(name string, channel chan<- model.Status, wg *sync.WaitGroup) {
	defer wg.Done()
	var state model.Status
	var err error
	var raw []byte

	raw, err = model.Request(name, message.TypeNodeStatus, []byte{})
	if err != nil {
		logger.Comm.Println("getNodeStatus", err)
		goto sendres
	}

	err = config.Jsoner.Unmarshal(raw, &state)
	if err != nil {
		logger.Exceptions.Println("UnmarshalNodeStatus", err)
		goto sendres
	}
sendres:
	channel <- state
}

func NodePrune(ctx *gin.Context) {
	model.PruneDeadNode()
	ctx.Status(200)
}
