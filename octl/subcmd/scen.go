package subcmd

import (
	"context"
	"strings"

	"github.com/piaodazhu/Octopoda/octl/node"
	"github.com/piaodazhu/Octopoda/octl/scenario"
)

func scenCmd(args []string) {
	var operation string
	if len(args) == 0 {
		goto usage
	}
	operation = args[0]
	args = args[1:]
	switch operation {
	case "apply":
		message, args := extractArgString(args, "-m", "--message", "")
		if len(args) == 0 {
			goto usage
		}

		deployment := args[0]
		if deployment != "purge" && message == "" {
			goto usage
		}

		target := "default"
		if len(args) > 2 {
			goto usage
		}
		if len(args) == 2 {
			target = args[1]
		}
		scenario.ScenarioApply(context.Background(), deployment, target, message)
	case "create":
		if len(args) == 0 {
			goto usage
		} else if len(args) == 1 {
			scenario.Create(args[0], "YourApp")
		} else if len(args) == 2 || args[1] != "with" {
			goto usage
		}
		scenario.Create(args[0], args[2:]...)
	case "get":
		if len(args) == 0 || args[0] == "ALL" {
			scenario.ScenariosInfo()
		} else if len(args) == 1 {
			scenario.ScenarioInfo(args[0])
		} else {
			goto usage
		}
	case "repo":
		user, args := extractArgString(args, "-u", "--username", "")
		if len(args) != 2 {
			goto usage
		}
		switch args[0] {
		case "clone":
			scenario.GitClone(args[1], user)
		case "push":
			scenario.GitPush(args[1], user)
		default:
			goto usage
		}
	case "version":
		if len(args) == 1 {
			scenario.ScenarioVersion(args[0])
		} else {
			goto usage
		}
	case "reset":
		message, args := extractArgString(args, "-m", "--message", "")
		version, args := extractArgString(args, "-v", "--version", "")
		if message == "" || version == "" {
			goto usage
		}
		if len(args) != 1 {
			goto usage
		}
		scenario.ScenarioReset(args[0], version, message)
	default:
		goto usage
	}
	return
usage:
	PrintUsage("scen")
}

func nappCmd(args []string) {
	var operation, nodename string
	if len(args) < 2 {
		goto usage
	}
	operation = args[0]
	args = args[1:]
	nodename = args[0]
	args = args[1:]
	switch operation {
	case "get":
		if len(args) == 0 || args[0] == "ALL" {
			node.NodeAppsInfo(nodename)
		} else if len(args) == 1 {
			nodeapp := args[0]
			appscen := strings.Split(nodeapp, "@")
			if len(appscen) == 2 && len(appscen[0]) != 0 && len(appscen[1]) != 0 {
				node.NodeAppInfo(nodename, appscen[0], appscen[1])
			} else {
				goto usage
			}
		} else {
			goto usage
		}
	// case "version":
	// 	if len(args) == 1 {
			
	// 	} else {
	// 		goto usage
	// 	}
	case "reset":
		message, args := extractArgString(args, "-m", "--message", "")
		version, args := extractArgString(args, "-v", "--version", "")
		if message == "" || version == "" {
			goto usage
		}
		if len(args) != 1 {
			goto usage
		}

		nodeapp := args[0]
		appscen := strings.Split(nodeapp, "@")
		if len(appscen) == 2 && len(appscen[0]) != 0 && len(appscen[1]) != 0 {
			node.NodeAppInfo(nodename, appscen[0], appscen[1])
		} else {
			goto usage
		}
		node.NodeAppReset(nodename, appscen[0], appscen[1], version, message)
	default:
		goto usage
	}
	return
usage:
	PrintUsage("napp")
}
