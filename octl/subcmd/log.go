package subcmd

import (
	"github.com/piaodazhu/Octopoda/octl/log"
)

func logCmd(args []string) {
	var class string = "Log"
	var operation string = ""
	if len(args) == 0 {
		goto usage
	}
	operation = args[0]
	args = args[1:]
	switch operation {
	case "get":
		maxlines, args := extractArgInt(args, "-l", "--lines", 30)
		maxdaysbefore, args := extractArgInt(args, "-d", "--days", 0)
		if len(args) != 1 {
			goto usage
		}
		log.PullLog(args[0], maxlines, maxdaysbefore)
	default:
		goto usage
	}
	return
usage:
	PrintUsage(class, operation)
}
