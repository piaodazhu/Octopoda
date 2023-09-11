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

func Create(arglist []string) {
	if len(arglist) == 0 {
		goto usage
	} else if len(arglist) == 1 {
		scenario.Create(arglist[0], "yourApp")
	} else if len(arglist) == 2 || arglist[1] != "with" {
		goto usage
	} else {
		scenario.Create(arglist[0], arglist[2:]...)
	}
	return
usage:
	PrintUsage("create")
}

func Apply(arglist []string) {
	if len(arglist) == 0 || len(arglist) > 4 {
		goto usage
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
usage:
	PrintUsage("apply")
}

func Get(arglist []string) {
	if len(arglist) == 0 {
		goto usage
	}
	switch arglist[0] {
	case "nodes":
		node.NodesInfo()
	case "node":
		if len(arglist) != 2 {
			goto usage
		}
		node.NodeInfo(arglist[1])
	case "scenarios":
		if len(arglist) != 1 {
			goto usage
		}
		scenario.ScenariosInfo()
	case "scenario":
		if len(arglist) != 2 {
			goto usage
		}
		scenario.ScenarioInfo(arglist[1])
	case "nodeapps":
		if len(arglist) != 2 {
			goto usage
		}
		node.NodeAppsInfo(arglist[1])
	case "nodeapp":
		if len(arglist) != 3 {
			goto usage
		}
		nodeapp := arglist[2]
		appscen := strings.Split(nodeapp, "@")
		if len(appscen) == 2 && len(appscen[0]) != 0 && len(appscen[1]) != 0 {
			node.NodeAppInfo(arglist[1], appscen[0], appscen[1])
		} else {
			goto usage
		}
	}
	return
usage:
	PrintUsage("get")
}

func Status(arglist []string) {
	if len(arglist) != 1 {
		goto usage
	}
	switch arglist[0] {
	case "nodes":
		node.NodesStatus()
	default:
		node.NodeStatus(arglist[0])
	}
	return
usage:
	PrintUsage("status")
}

func Fix(arglist []string) {
	if len(arglist) != 2 {
		goto usage
	}
	if arglist[0] == "scenario" {
		scenario.ScenarioFix(arglist[1])
	} else if arglist[0] == "node" {
		// node.NodeAppsFix(arglist[1])
		fmt.Println("not support node fix")
	}
	return
usage:
	PrintUsage("fix")
}

func Log(arglist []string) {
	if len(arglist) == 0 {
		goto usage
	}
	log.NodeLog(arglist[0], arglist[1:])
	return
usage:
	PrintUsage("log")
}

func Version(arglist []string) {
	if len(arglist) < 2 {
		goto usage
	}
	if arglist[0] == "scenario" && len(arglist) == 2 {
		scenario.ScenarioVersion(arglist[1])
	} else if arglist[0] == "nodeapp" && len(arglist) == 3 {
		nodeapp := arglist[2]
		appscen := strings.Split(nodeapp, "@")
		if len(appscen) == 2 && len(appscen[0]) != 0 && len(appscen[1]) != 0 {
			node.NodeAppInfo(arglist[1], appscen[0], appscen[1])
		} else {
			goto usage
		}
	}
	return
usage:
	PrintUsage("version")
}

func Reset(arglist []string) {
	var message, version string
	if len(arglist) < 6 {
		goto usage
	}

	if arglist[len(arglist)-2] == "-m" {
		message = arglist[len(arglist)-1]
	} else {
		goto usage
	}

	if arglist[len(arglist)-4] == "-v" {
		version = arglist[len(arglist)-3]
	} else {
		goto usage
	}

	if arglist[0] == "scenario" && len(arglist) == 6 {
		scenario.ScenarioReset(arglist[1], version, message)
	} else if arglist[0] == "nodeapp" && len(arglist) == 7 {
		nodeapp := arglist[2]
		appscen := strings.Split(nodeapp, "@")
		if len(appscen) == 2 && len(appscen[0]) != 0 && len(appscen[1]) != 0 {
			node.NodeAppReset(arglist[1], appscen[0], appscen[1], version, message)
		} else {
			goto usage
		}
	}
	return
usage:
	PrintUsage("reset")
}

func SetSSH(arglist []string) {
	if len(arglist) != 1 {
		goto usage
	}
	shell.SetSSH(arglist[0])
	return
usage:
	PrintUsage("setssh")
}

func GetSSH(arglist []string) {
	if len(arglist) != 0 {
		goto usage
	}
	shell.GetSSH()
	return
usage:
	PrintUsage("ssh")
}

func DelSSH(arglist []string) {
	if len(arglist) != 1 {
		goto usage
	}
	shell.DelSSH(arglist[0])
	return
usage:
	PrintUsage("setssh")
}

func SSH(arglist []string) {
	if len(arglist) != 1 {
		goto usage
	}
	shell.SSH(arglist[0])
	return
usage:
	PrintUsage("ssh")
}

func Upload(arglist []string) {
	if len(arglist) != 2 {
		goto usage
	}
	file.UpLoadFile(arglist[0], arglist[1])
	return
usage:
	PrintUsage("upload")
}

func Spread(arglist []string) {
	if len(arglist) < 3 {
		goto usage
	}
	file.SpreadFile(arglist[0], arglist[1], arglist[2:])
	return
usage:
	PrintUsage("spread")
}

func Distrib(arglist []string) {
	if len(arglist) < 3 {
		goto usage
	}
	file.DistribFile(arglist[0], arglist[1], arglist[2:])
	return
usage:
	PrintUsage("distrib")
}

func FileTree(arglist []string) {
	if len(arglist) < 2 {
		goto usage
	}
	switch arglist[0] {
	case "store":
		if len(arglist) == 2 {
			file.ListAllFile(arglist[0], arglist[1], "")
		} else if len(arglist) == 3 {
			file.ListAllFile(arglist[0], arglist[1], arglist[2])
		} else {
			goto usage
		}
	case "log":
		if len(arglist) == 2 {
			file.ListAllFile(arglist[0], arglist[1], "")
		} else {
			goto usage
		}
	case "nodeapp":
		if len(arglist) == 3 {
			file.ListAllFile(arglist[0], arglist[1], arglist[2])
		} else if len(arglist) == 4 {
			file.ListAllFile(arglist[0], arglist[1], arglist[2]+"/"+arglist[3])
		} else {
			goto usage
		}
	default:
		goto usage
	}
	return
usage:
	PrintUsage("tree")
}

func Pull(arglist []string) {
	if len(arglist) < 3 {
		goto usage
	}
	switch arglist[0] {
	case "store":
		fallthrough
	case "log":
		if len(arglist) == 4 {
			file.PullFile(arglist[0], arglist[1], arglist[2], arglist[3])
		} else {
			goto usage
		}
	case "nodeapp":
		if len(arglist) == 4 {
			file.PullFile(arglist[0], arglist[1], arglist[2], arglist[3])
		} else if len(arglist) == 5 {
			file.PullFile(arglist[0], arglist[1], arglist[2]+"/"+arglist[3], arglist[4])
		} else {
			goto usage
		}
	default:
		goto usage
	}
	return
usage:
	PrintUsage("pull")
}

func Prune(arglist []string) {
	if len(arglist) != 0 {
		goto usage
	}
	node.NodePrune()
	return
usage:
	PrintUsage("prune")
}

func Run(arglist []string) {
	if len(arglist) < 2 {
		goto usage
	}
	shell.Run(arglist[0], arglist[1:])
	return
usage:
	PrintUsage("run")
}

func XRun(arglist []string) {
	if len(arglist) < 2 {
		goto usage
	}
	shell.XRun(arglist[0], arglist[1:])
	return
usage:
	PrintUsage("run")
}


func Pakma(arglist []string) {
	if len(arglist) < 2 {
		goto usage
	}
	node.Pakma(arglist[0], arglist[1:])
	return
usage:
	PrintUsage("pakma")
}

func Group(arglist []string) {
	switch arglist[0] {
	case "get":
		if len(arglist) != 2 {
			goto usage
		}
		node.GroupGet(arglist[1])
	case "get-all":
		if len(arglist) != 1 {
			goto usage
		}
		node.GroupGetAll()
	case "set":
		if len(arglist) <= 2 {
			goto usage
		}
		node.GroupSet(arglist[1], false, arglist[2:])
	case "set-nocheck":
		if len(arglist) <= 2 {
			goto usage
		}
		node.GroupSet(arglist[1], true, arglist[2:])
	case "del":
		if len(arglist) != 2 {
			goto usage
		}
		node.GroupDel(arglist[1])
	default:
		goto usage
	}
	return
usage:
	PrintUsage("group")
}
