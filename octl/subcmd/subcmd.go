package subcmd

import (
	"fmt"
	"octl/file"
	"octl/log"
	"octl/node"
	"octl/scenario"
	"octl/shell"
	"strings"
)

func Apply(arglist []string) {
	if len(arglist) == 0 || len(arglist) > 4 {
		return
	}
	if arglist[len(arglist)-2] == "-m" {
		if len(arglist) == 3 {
			scenario.ScenarioApply(arglist[0], "default", arglist[2])
			return
		} else if len(arglist) == 4 {
			scenario.ScenarioApply(arglist[0], arglist[1], arglist[3])
			return
		}
	} else if arglist[1] == "purge" {
		scenario.ScenarioApply(arglist[0], arglist[1], "")
		return 
	}
	fmt.Println(`usage: octl apply xx.yaml [target] -m "your message"`)
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
	case "scenarios":
		if len(arglist) != 1 {
			return
		}
		scenario.ScenariosInfo()
	case "scenario":
		if len(arglist) != 2 {
			return
		}
		scenario.ScenarioInfo(arglist[1])
	case "nodeapps":
		if len(arglist) != 2 {
			return
		}
		node.NodeAppsInfo(arglist[1])
	case "nodeapp":
		if len(arglist) != 3 {
			return
		}
		nodeapp := arglist[2]
		appscen := strings.Split(nodeapp, "@")
		if len(appscen) == 2 && len(appscen[0]) != 0 && len(appscen[1]) != 0 {
			node.NodeAppInfo(arglist[1], appscen[0], appscen[1])
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

func Fix(arglist []string) {
	if len(arglist) != 2 {
		return
	}
	if arglist[0] == "scenario" {
		scenario.ScenarioFix(arglist[1])
	} else if arglist[0] == "node" {
		// node.NodeAppsFix(arglist[1])
		fmt.Println("not support node fix")
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
	// case "scenario":

	// default:
	}
}

func Version(arglist []string) {
	if len(arglist) < 2 {
		return
	}
	if arglist[0] == "scenario" && len(arglist) == 2 {
		scenario.ScenarioVersion(arglist[1])
	} else if arglist[0] == "nodeapp" && len(arglist) == 3 {
		nodeapp := arglist[2]
		appscen := strings.Split(nodeapp, "@")
		if len(appscen) == 2 && len(appscen[0]) != 0 && len(appscen[1]) != 0 {
			node.NodeAppInfo(arglist[1], appscen[0], appscen[1])
			return
		}
	}
}

func Reset(arglist []string) {
	var message, version string 
	if len(arglist) < 6 {
		goto printusage
	}
	
	if arglist[len(arglist) - 2] == "-m" {
		message = arglist[len(arglist) - 1]
	} else {
		goto printusage
	}

	if arglist[len(arglist) - 4] == "-v" {
		version = arglist[len(arglist) - 3]
	} else {
		goto printusage
	}

	if arglist[0] == "scenario" && len(arglist) == 6 {
		scenario.ScenarioReset(arglist[1], version, message)
		return
	} else if arglist[0] == "nodeapp" && len(arglist) == 7 {
		nodeapp := arglist[2]
		appscen := strings.Split(nodeapp, "@")
		if len(appscen) == 2 && len(appscen[0]) != 0 && len(appscen[1]) != 0 {
			node.NodeAppReset(arglist[1], appscen[0], appscen[1], version, message)
			return
		}
	}
printusage:
	fmt.Println("Usage: octl reset [scenario <scen>|nodeapp <node> <app>@<scen>]  -v <version> -m <message>")
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

func Distrib(arglist []string) {
	if len(arglist) < 3 {
		return
	}
	file.DistribFile(arglist[0], arglist[1], arglist[2:])
}

func FileTree(arglist []string) {
	if len(arglist) == 0 {
		return
	} else if len(arglist) == 1 {
		file.ListAllFile(arglist[0], "")
	} else if len(arglist) == 2 {
		file.ListAllFile(arglist[0], arglist[1])
	} else {
		return
	}
}

func Prune(arglist []string) {
	if len(arglist) != 0 {
		return
	}
	node.NodePrune()
}

func Run(arglist []string) {
	if len(arglist) < 2 {
		return
	}
	shell.RunTask(arglist[0], arglist[1:])
}
