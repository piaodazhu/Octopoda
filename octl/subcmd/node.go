package subcmd

import "github.com/piaodazhu/Octopoda/octl/node"

func nodeCmd(args []string) {
	var operation string
	if len(args) == 0 {
		goto usage
	}
	operation = args[0]
	args = args[1:]
	switch operation {
	case "get":
		hf, args := extractArgString(args, "-hf", "--health-filter", "")
		mf, args := extractArgString(args, "-mf", "--msgconn-filter", "")
		if len(args) == 0 || args[0] == "ALL" {
			node.NodesInfoWithFilter(nil, hf, mf)
		} else {
			node.NodesInfoWithFilter(args, hf, mf)
		}
	case "prune":
		if len(args) == 0 {
			goto usage
		}
		if args[0] == "ALL" {
			node.NodesPrune(nil)
		} else {
			node.NodesPrune(args)
		}
	case "status":
		if len(args) == 0 || args[0] == "ALL" {
			node.NodesStatus(nil)
		} else {
			node.NodesStatus(args)
		}
	default:
		goto usage
	}
	return
usage:
	PrintUsage("node")
}
