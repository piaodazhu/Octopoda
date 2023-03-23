package model

import (
	"brain/config"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

var NodeMap map[string]*NodeModel
var Lock sync.RWMutex

func InitNodeMap() {
	NodeMap = make(map[string]*NodeModel)
	go func() {
		for {
			time.Sleep(time.Second)
			Lock.Lock()
			for _, node := range NodeMap {
				if node.State == NodeStateDisconn && node.OfflineTs+int64(config.GlobalConfig.TentacleFace.RecordTimeout) < time.Now().Unix() {
					// logger.Tentacle.Print("MarkDeadNode", nodename)
					node.State = NodeStateDead
				}
			}
			Lock.Unlock()
		}
	}()
}

func StoreNode(name string, ip string, port uint16) uint32 {
	var node *NodeModel
	var sb strings.Builder
	sb.WriteString(ip)
	sb.WriteByte(':')
	sb.WriteString(strconv.Itoa(int(port)))

	Lock.Lock()
	defer Lock.Unlock()

	if n, found := NodeMap[name]; found {
		node = n
	} else {
		node = &NodeModel{
			Id:      uuid.New().ID(),
			Name:    name,
			Applist: []*AppModel{},
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
	Lock.Lock()
	defer Lock.Unlock()
	if node, found := NodeMap[name]; found {
		node.State = int32(NodeStateDisconn)
		node.OfflineTs = time.Now().Unix() //
		return true
	}
	return false
}

func PruneDeadNode() {
	Lock.RLock()
	defer Lock.RUnlock()

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
	Lock.RLock()
	defer Lock.RUnlock()

	if node, found := NodeMap[name]; found {
		return node, true
	}
	return nil, false
}

func GetNodeInfoById(id int) (*NodeModel, bool) {
	Lock.RLock()
	defer Lock.RUnlock()

	for _, node := range NodeMap {
		if node.Id == uint32(id) {
			return node, true
		}
	}
	return nil, false
}

func GetNodesInfoAll() ([]*NodeModel, bool) {
	Lock.RLock()
	defer Lock.RUnlock()

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

func GetNodeAddress(name string) (string, bool) {
	Lock.Lock()
	defer Lock.Unlock()

	if node, found := NodeMap[name]; found {
		return node.Addr, true
	}
	return "", false
}

func GetNodeState(name string) (int, bool) {
	Lock.RLock()
	defer Lock.RUnlock()

	if node, found := NodeMap[name]; found {
		return int(node.State), true
	}
	return -1, false
}

func SetNodeState(name string, state int) bool {
	Lock.Lock()
	defer Lock.Unlock()

	if node, found := NodeMap[name]; found {
		node.State = int32(state)
		return true
	}
	return false
}
