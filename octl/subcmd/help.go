package subcmd

import "fmt"

var usageInfo = map[string]string{
	"apply":   `usage: octl apply xx.yaml [target] -m "your message"`,
	"get":     `usage: octl get [nodes|node <node>|scenarios|scenario <scen>|nodeapps <node>|nodeapp <node> <app>@<scen>]"`,
	"status":  `usage: octl status [nodes|node <node>]`,
	"fix":     `usage: octl fix scenario <scen>`,
	"log":     `usage: octl log [master|node <node>] l<maxline> d<maxday>`,
	"version": `usage: octl version [scenario <scen>|nodeapp <node> <app>@<scen>]`,
	"reset":   `usage: octl reset [scenario <scen>|nodeapp <node> <app>@<scen>]  -v <version> -m <message>`,
	"shell":   `usage: octl shell <node>`,
	"upload":  `usage: octl upload <localFileOrDir> <targetDir>`,
	"spread":  `usage: octl spread <masterFileOrDir> <targetDir> <node1> <node2> ...`,
	"distrib": `usage: octl distrib <localFileOrDir> <targetDir> <node1> <node2> ...`,
	"tree":    `usage: octl tree [master|<node>] [SubDir]`,
	"prune":   `usage: octl prune`,
	"run":     `usage: octl run [ {<command>} | <script> ] <node1> <node2> ...`,
}

func PrintUsage(subcmd string) {
	fmt.Println(usageInfo[subcmd])
}

func PrintUsages() {
	fmt.Println("Subcmd:")
	for subcmd, usage := range usageInfo {
		fmt.Printf("\t%s:\n\t\t%s\n", subcmd, usage)
	}
}
