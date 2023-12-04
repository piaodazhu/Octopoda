package subcmd

import "github.com/piaodazhu/Octopoda/octl/node"

// TODO refactor
func pakmaCmd(args []string) {
	var class string = "PAcKage MAnager"
	var operation string = ""
	if len(args) == 0 {
		goto usage
	}
	operation = args[0]
	args = args[1:]
	node.Pakma(operation, args)
	return
usage:
	PrintUsage(class, operation)
}
