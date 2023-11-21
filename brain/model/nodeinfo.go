package model

import (
	"net"
	"sync"
	"time"

	"github.com/piaodazhu/Octopoda/brain/config"
	"github.com/piaodazhu/Octopoda/protocols"
)

type NodeModel struct {
	protocols.NodeInfo
	// MsgConn *net.Conn
	ConnInfo
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
				if node.State == protocols.NodeStateDisconn && node.OfflineTs+int64(config.GlobalConfig.TentacleFace.RecordTimeout) < time.Now().UnixMilli() {
					// logger.Tentacle.Print("MarkDeadNode", nodename)
					node.State = protocols.NodeStateDead
					node.ConnInfo.Close()
				}
			}
			NodesLock.Unlock()
		}
	}()
}

func StoreNode(name, version, addr string, conn net.Conn) {
	var node *NodeModel

	NodesLock.Lock()
	defer NodesLock.Unlock()

	if n, found := NodeMap[name]; found {
		node = n
		if conn != nil {
			node.ConnInfo.Fresh(conn)
			node.ConnInfo.StartReceive()
		}
	} else {
		info := protocols.NodeInfo{
			Name: name,
		}
		node = &NodeModel{
			NodeInfo: info,
			// Applist:   []*AppModel{},
		}
		node.ConnInfo = CreateConnInfo(conn, &node.ConnState)
		NodeMap[name] = node
	}

	if conn == nil { // 只对于heartbeat通道做状态设置
		node.Version = version
		node.Addr = addr
		node.State = protocols.NodeStateDisconn // 必须心跳后才能认为上线
		node.OnlineTs = time.Now().UnixMilli()
	}
}

func UpdateNode(name string, delay int64) bool {
	// logger.Tentacle.Print("UpdateNode", name)
	return SetNodeStateAndDelay(name, protocols.NodeStateReady, delay)
}

func DisconnNode(name string) bool {
	// logger.Tentacle.Print("DisconnNode", name)
	NodesLock.Lock()
	defer NodesLock.Unlock()
	if node, found := NodeMap[name]; found {
		node.State = int32(protocols.NodeStateDisconn)
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
		if node.State == protocols.NodeStateDead {
			toBePruned = append(toBePruned, name)
		}
	}
	for i := range toBePruned {
		delete(NodeMap, toBePruned[i])
	}
}

func GetNodeInfoByName(name string) (protocols.NodeInfo, bool) {
	NodesLock.RLock()
	defer NodesLock.RUnlock()

	if node, found := NodeMap[name]; found {
		return node.NodeInfo, true
	}
	return protocols.NodeInfo{}, false
}

func GetNodesInfoAll() ([]protocols.NodeInfo, bool) {
	NodesLock.RLock()
	defer NodesLock.RUnlock()

	if len(NodeMap) == 0 {
		return nil, false
	}
	res := make([]protocols.NodeInfo, 0, len(NodeMap))
	for _, val := range NodeMap {
		copynode := val.NodeInfo
		res = append(res, copynode)
	}
	return res, true
}

func GetNodesInfo(names []string) ([]protocols.NodeInfo, bool) {
	NodesLock.RLock()
	defer NodesLock.RUnlock()

	if len(NodeMap) == 0 {
		return nil, false
	}
	res := make([]protocols.NodeInfo, 0, len(names))
	for _, name := range names {
		if node, found := NodeMap[name]; found {
			copynode := node.NodeInfo
			res = append(res, copynode)
		}
	}
	return res, true
}

var (
	GetConnOk     = 0
	GetConnNoNode = 1
	GetConnNoConn = 2
)

func GetNodeMsgConn(name string) (*ConnInfo, int) {
	NodesLock.RLock()
	defer NodesLock.RUnlock()
	node, found := NodeMap[name]
	if !found {
		return nil, GetConnNoNode
	}
	if node.ConnState == "Off" {
		return &node.ConnInfo, GetConnNoConn
	}
	return &node.ConnInfo, GetConnOk
}

func ResetNodeMsgConn(name string) {
	NodesLock.RLock()
	defer NodesLock.RUnlock()
	node, found := NodeMap[name]
	if found {
		node.Close()
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

func SetNodeStateAndDelay(name string, state int, delay int64) bool {
	NodesLock.Lock()
	defer NodesLock.Unlock()

	if node, found := NodeMap[name]; found {
		node.State = int32(state)
		node.ActiveTs = time.Now().UnixMilli()
		node.Delay = delay
		return true
	}
	return false
}

func GetNodesMaxDelay(names []string) int64 {
	NodesLock.Lock()
	defer NodesLock.Unlock()

	var maxDelay int64 = 0
	for _, name := range names {
		if node, found := NodeMap[name]; found {
			if node.State == protocols.NodeStateReady && node.Delay > maxDelay {
				maxDelay = node.Delay
			}
		}
	}
	return maxDelay
}