package subcmd

import (
	"github.com/piaodazhu/Octopoda/octl/workgroup"
)

func workgroupCmd(args []string) {
	var class string = "Workgroup"
	var operation, path string = "", ""
	if len(args) == 0 {
		goto usage
	}
	operation = args[0]
	args = args[1:]
	switch operation {
	case "pwd":
		workgroup.Pwd()
	case "cd":
		if len(args) == 0 {
			path = "~"
		} else if len(args) == 1 {
			path = args[0]
		} else {
			goto usage
		}
		workgroup.Cd(path)
	case "ls":
		if len(args) == 0 {
			path = ""
		} else if len(args) == 1 {
			path = args[0]
		} else {
			goto usage
		}
		workgroup.Ls(path)
	case "get":
		if len(args) == 0 {
			path = ""
		} else if len(args) == 1 {
			path = args[0]
		} else {
			goto usage
		}
		workgroup.Get(path)
	case "add":
		if len(args) > 1 {
			path = args[0]
			workgroup.AddMembers(path, args[1:]...)
		} else {
			goto usage
		}
	case "rm":
		if len(args) >= 1 {
			path = args[0]
			// TODO: test if nil
			workgroup.RemoveMembers(path, args[1:]...)
		} else {
			goto usage
		}
	case "grant":
		if len(args) == 2 {
			path = args[0]
			passwd := args[1]
			workgroup.Grant(path, passwd)
		} else {
			goto usage
		}
	default:
		goto usage
	}
	return
usage:
	PrintUsage(class, operation)
}
