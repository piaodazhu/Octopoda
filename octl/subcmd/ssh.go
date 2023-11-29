package subcmd

import "github.com/piaodazhu/Octopoda/octl/shell"

func sshCmd(args []string) {
	var class string = "SSH"
	var operation string = ""
	if len(args) == 0 {
		goto usage
	}
	operation = args[0]
	args = args[1:]
	switch operation {
	case "login":
		if len(args) != 1 {
			goto usage
		}
		shell.SSH(args[0])
	case "set":
		if len(args) != 1 {
			goto usage
		}
		shell.SetSSH(args[0])
	case "ls":
		if len(args) != 0 {
			goto usage
		}
		shell.GetSSH()
	case "del":
		if len(args) != 1 {
			goto usage
		}
		shell.DelSSH(args[0])
	default:
		goto usage
	}
	return
usage:
	PrintUsage(class, operation)
}
