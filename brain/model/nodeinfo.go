package model

import (
	"brain/config"
	"net"
	"sync"
	"time"
)

const (
	NodeStateReady = iota
	NodeStateDisconn
	NodeStateDead
)

type NodeInfo struct {
	Name      string
	State     int32
	OnlineTs  int64
	OfflineTs int64
	ActiveTs  int64
}

type NodeModel struct {
	NodeInfo
	MsgConn *net.Conn
}

type State struct {
	Name      string
	Platform  string
	CpuCores  int
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
				if node.State == NodeStateDisconn && node.OfflineTs+int64(config.GlobalConfig.TentacleFace.RecordTimeout) < time.Now().UnixMilli() {
					// logger.Tentacle.Print("MarkDeadNode", nodename)
					node.State = NodeStateDead
				}
			}
			NodesLock.Unlock()
		}
	}()
}

func StoreNode(name string, conn *net.Conn) {
	var node *NodeModel

	NodesLock.Lock()
	defer NodesLock.Unlock()

	if n, found := NodeMap[name]; found {
		node = n
		if n.MsgConn != nil {
			(*n.MsgConn).Close()
		}
	} else {
		info := NodeInfo{
			Name: name,
		}
		node = &NodeModel{
			NodeInfo: info,
			// Applist:   []*AppModel{},
		}
		NodeMap[name] = node
	}
	node.MsgConn = conn
	node.State = NodeStateReady
	node.OnlineTs = time.Now().UnixMilli()
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
		node.OfflineTs = time.Now().UnixMilli() //
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

func GetNodesInfoAll() ([]*NodeModel, bool) {
	NodesLock.RLock()
	defer NodesLock.RUnlock()

	if len(NodeMap) == 0 {
		return nil, false
	}
	res := make([]*NodeModel, 0, len(NodeMap))
	for _, val := range NodeMap {
		copynode := *val
		res = append(res, &copynode)
	}
	return res, true
}

var (
	GetConnOk = 0
	GetConnNoNode = 1
	GetConnNoConn = 2
)
func GetNodeMsgConn(name string) (*net.Conn, int) {
	NodesLock.RLock()
	defer NodesLock.RUnlock()
	node, found := NodeMap[name]
	if !found {
		return nil, GetConnNoNode
	}
	if node.MsgConn == nil {
		return nil, GetConnNoConn
	}
	return node.MsgConn, GetConnOk
}

func ResetNodeMsgConn(name string) {
	NodesLock.RLock()
	defer NodesLock.RUnlock()
	node, found := NodeMap[name]
	if found && node.MsgConn != nil {
		(*node.MsgConn).Close()
		node.MsgConn = nil
	}
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
		node.ActiveTs = time.Now().UnixMilli()
		return true
	}
	return false
}
