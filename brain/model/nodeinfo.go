package model

import (
	"brain/config"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	NodeStateReady = iota
	NodeStateDisconn
	NodeStateDead
)

type NodeInfo struct {
	Id        uint32
	Name      string
	Addr      string
	State     int32
	OnlineTs  int64
	OfflineTs int64
}

type NodeModel struct {
	NodeInfo
	// Applist []*AppModel
}

type State struct {
	Id        int
	Name      string
	Platform  string
	CpuCores  int
	Ip        string
	LocalTime int64

	CpuLoadShort float64
	CpuLoadLong  float64
	MemUsed      uint64
	MemTotal     uint64
	DiskUsed     uint64
	DiskTotal    uint64
}


var NodeMap map[string]*NodeModel
var NodesLock sync.RWMutex

func InitNodeMap() {
	NodeMap = make(map[string]*NodeModel)
	go func() {
		for {
			time.Sleep(time.Second)
			NodesLock.Lock()
			for _, node := range NodeMap {
				if node.State == NodeStateDisconn && node.OfflineTs+int64(config.GlobalConfig.TentacleFace.RecordTimeout) < time.Now().Unix() {
					// logger.Tentacle.Print("MarkDeadNode", nodename)
					node.State = NodeStateDead
				}
			}
			NodesLock.Unlock()
		}
	}()
}

func StoreNode(name string, ip string, port uint16) uint32 {
	var node *NodeModel
	var sb strings.Builder
	sb.WriteString(ip)
	sb.WriteByte(':')
	sb.WriteString(strconv.Itoa(int(port)))

	NodesLock.Lock()
	defer NodesLock.Unlock()

	if n, found := NodeMap[name]; found {
		node = n
	} else {
		info := NodeInfo{
			Id:   uuid.New().ID(),
			Name: name,
		}
		node = &NodeModel{
			NodeInfo: info,
			// Applist:   []*AppModel{},
		}
		NodeMap[name] = node
	}
	node.State = NodeStateReady
	node.Addr = sb.String()
	node.OnlineTs = time.Now().Unix()

	return node.Id
}

func UpdateNode(name string) bool {
	// logger.Tentacle.Print("UpdateNode", name)
	return SetNodeState(name, NodeStateReady)
}

func DisconnNode(name string) bool {
	// logger.Tentacle.Print("DisconnNode", name)
	NodesLock.Lock()
	defer NodesLock.Unlock()
	if node, found := NodeMap[name]; found {
		node.State = int32(NodeStateDisconn)
		node.OfflineTs = time.Now().Unix() //
		return true
	}
	return false
}

func PruneDeadNode() {
	NodesLock.RLock()
	defer NodesLock.RUnlock()

	toBePruned := []string{}
	for name, node := range NodeMap {
		if node.State == NodeStateDead {
			toBePruned = append(toBePruned, name)
		}
	}
	for i := range toBePruned {
		delete(NodeMap, toBePruned[i])
	}
}

func GetNodeInfoByName(name string) (*NodeModel, bool) {
	NodesLock.RLock()
	defer NodesLock.RUnlock()

	if node, found := NodeMap[name]; found {
		return node, true
	}
	return nil, false
}

func GetNodeInfoById(id int) (*NodeModel, bool) {
	NodesLock.RLock()
	defer NodesLock.RUnlock()

	for _, node := range NodeMap {
		if node.Id == uint32(id) {
			return node, true
		}
	}
	return nil, false
}

func GetNodesInfoAll() ([]*NodeInfo, bool) {
	NodesLock.RLock()
	defer NodesLock.RUnlock()

	if len(NodeMap) == 0 {
		return nil, false
	}
	res := make([]*NodeInfo, 0, len(NodeMap))
	for _, val := range NodeMap {
		copynode := *val
		res = append(res, &copynode.NodeInfo)
	}
	return res, true
}

func GetNodeAddress(name string) (string, bool) {
	NodesLock.Lock()
	defer NodesLock.Unlock()

	if node, found := NodeMap[name]; found && node.State == NodeStateReady {
		return node.Addr, true
	}
	return "", false
}

func GetNodeState(name string) (int, bool) {
	NodesLock.RLock()
	defer NodesLock.RUnlock()

	if node, found := NodeMap[name]; found {
		return int(node.State), true
	}
	return -1, false
}

func SetNodeState(name string, state int) bool {
	NodesLock.Lock()
	defer NodesLock.Unlock()

	if node, found := NodeMap[name]; found {
		node.State = int32(state)
		return true
	}
	return false
}
