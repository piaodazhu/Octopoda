package subcmd

import "github.com/piaodazhu/Octopoda/octl/file"

func fileCmd(args []string) {
	var class string = "File"
	var operation string = ""
	if len(args) == 0 {
		goto usage
	}
	operation = args[0]
	args = args[1:]
	switch operation {
	case "upload":
		isForce, args := extractArgBool(args, "-f", "--force", true)
		if len(args) < 3 {
			goto usage
		}
		file.Upload(args[0], args[1], args[2:], isForce)
	case "download":
		if len(args) < 3 {
			goto usage
		}
		file.Download(args[0], args[1], args[2])
	case "ls":
		if len(args) != 2 {
			goto usage
		}
		file.ListAllFile(args[0], args[1])
	default:
		goto usage
	}
	return
usage:
	PrintUsage(class, operation)
}
