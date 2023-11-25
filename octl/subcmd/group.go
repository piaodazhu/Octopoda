package subcmd

import "github.com/piaodazhu/Octopoda/octl/node"

func groupCmd(args []string) {
	var operation string
	if len(args) == 0 {
		goto usage
	}
	operation = args[0]
	args = args[1:]
	switch operation {
	case "get":
		if len(args) == 0 || (len(args) == 1 && args[0] == "ALL") {
			node.GroupGetAll()
		} else if len(args) == 1 {
			node.GroupGet(args[0])
		} else {
			goto usage
		}
	case "set":
		nocheck, args := extractArgBool(args, "-n", "--no-check", false)
		if len(args) < 2 {
			goto usage
		}
		node.GroupSet(args[0], nocheck, args[1:])
	case "del":
		if len(args) == 1 {
			node.GroupDel(args[0])
		} else {
			goto usage
		}
	default:
		goto usage
	}
	return
usage:
	PrintUsage("group")
}
