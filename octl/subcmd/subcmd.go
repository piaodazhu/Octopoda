package subcmd

import (
	"github.com/piaodazhu/Octopoda/octl/output"
)

func Execute(subcmd string, args []string) {
	switch subcmd {
	case "node":
		nodeCmd(args)
	case "wg":
		workgroupCmd(args)
	case "cmd":
		cmdCmd(args)
	case "file":
		fileCmd(args)
	case "ssh":
		sshCmd(args)
	case "scen":
		scenCmd(args)
	case "napp":
		nappCmd(args)
	case "log":
		logCmd(args)
	case "pakma":
		pakmaCmd(args)
	case "help":
		PrintUsages(args)
	default:
		output.PrintFatalf("subcommand %s not support\n", subcmd)
	}
}
