package subcmd

import (
	"github.com/piaodazhu/Octopoda/octl/shell"
)

func cmdCmd(args []string) {
	var class string = "Command"
	var operation string = ""
	var cmd string
	if len(args) == 0 {
		goto usage
	}
	operation = args[0]
	args = args[1:]
	switch operation {
	case "run":
		timeAlign, args := extractArgBool(args, "-ta", "--time-align", false)
		if len(args) < 2 {
			goto usage
		}
		cmd, args = extractArgString(args, "-cc", "--common-command", "")
		if cmd != "" {
			shell.RunCommand(cmd, false, timeAlign, args)
			return
		}
		cmd, args = extractArgString(args, "-bg", "--background", "")
		if cmd != "" {
			shell.RunCommand(cmd, true, timeAlign, args)
			return
		}
		cmd, args = extractArgString(args, "-ss", "--shellscript", "")
		if cmd != "" {
			shell.RunScript(cmd, timeAlign, args)
			return
		}
		shell.RunCommand(args[0], false, timeAlign, args[1:])
	case "xrun":
		timeAlign, args := extractArgBool(args, "-ta", "--time-align", false)
		delayExec, args := extractArgInt(args, "-d", "--delay-execute", 0)
		if len(args) < 2 {
			goto usage
		}
		cmd, args = extractArgString(args, "-cc", "--common-command", "")
		if cmd != "" {
			shell.XRunCommand(cmd, false, timeAlign, delayExec, args)
			return
		}
		cmd, args = extractArgString(args, "-bg", "--background", "")
		if cmd != "" {
			shell.XRunCommand(cmd, true, timeAlign, delayExec, args)
			return
		}
		cmd, args = extractArgString(args, "-ss", "--shellscript", "")
		if cmd != "" {
			shell.XRunScript(cmd, timeAlign, delayExec, args)
			return
		}
		shell.XRunCommand(args[0], false, timeAlign, delayExec, args[1:])
	default:
		goto usage
	}
	return
usage:
	PrintUsage(class, operation)
}
