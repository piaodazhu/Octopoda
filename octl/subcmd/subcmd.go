package subcmd

import (
	"octl/file"
	"octl/log"
	"octl/node"
	"octl/shell"
)

func Apply(arglist []string) {

}

func Get(arglist []string) {
	if len(arglist) == 0 {
		return
	}
	switch arglist[0] {
	case "nodes":
		node.NodesInfo()
	case "node":
		if len(arglist) != 2 {
			return
		}
		node.NodeInfo(arglist[1])
	case "scienarios":
	case "scienario":
		if len(arglist) != 2 {
			return
		}
	}
}

func Status(arglist []string) {
	if len(arglist) == 0 {
		return
	}
	switch arglist[0] {
	case "nodes":
		node.NodesStatus()
	case "node":
		if len(arglist) != 2 {
			return
		}
		node.NodeStatus(arglist[1])
	}
}

func Log(arglist []string) {
	if len(arglist) == 0 {
		return
	}
	switch arglist[0] {
	case "master":
		log.NodeLog(arglist[0], arglist[1:])
	case "node":
		if len(arglist) == 1 {
			return
		}
		log.NodeLog(arglist[1], arglist[2:])
	case "scenario":

	default:

	}
}

func Version(arglist []string) {

}

func Reset(arglist []string) {

}

func Shell(arglist []string) {
	if len(arglist) != 1 {
		return
	}
	shell.SSH(arglist[0])
}

func Upload(arglist []string) {
	if len(arglist) != 2 {
		return
	}
	file.UpLoadFile(arglist[0], arglist[1])
}

func Spread(arglist []string) {
	if len(arglist) < 4 {
		return
	}
	file.SpreadFile(arglist[0], arglist[1], arglist[2], arglist[3:])
}

func Prune(arglist []string) {
	if len(arglist) != 0 {
		return
	}
	node.NodePrune()
}
