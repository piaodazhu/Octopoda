package subcmd

import "fmt"

var usageInfo = map[string]string{
	"apply":   `octl apply xx.yaml [target] -m "your message"`,
	"get":     `octl get [nodes|node <node>|scenarios|scenario <scen>|nodeapps <node>|nodeapp <node> <app>@<scen>]`,
	"status":  `octl status [nodes|node <node>]`,
	"fix":     `octl fix scenario <scen>`,
	"log":     `octl log [master|node <node>] [l<maxline>] [d<maxday>]`,
	"version": `octl version [scenario <scen>|nodeapp <node> <app>@<scen>]`,
	"reset":   `octl reset [scenario <scen>|nodeapp <node> <app>@<scen>]  -v <version> -m <message>`,
	"shell":   `octl shell <node>`,
	"upload":  `octl upload <localFileOrDir> <targetDir>`,
	"spread":  `octl spread <masterFileOrDir> <targetDir> <node1> <node2> ...`,
	"distrib": `octl distrib <localFileOrDir> <targetDir> <node1> <node2> ...`,
	"tree":    `octl tree [store [master|<node>]|log [<node>|master]|nodeapp <node> <app>@<scen>] [SubDir]`,
	"prune":   `octl prune`,
	"run":     `octl run [ {<command>} | <script> ] <node1> <node2> ...`,
}

func PrintUsage(subcmd string) {
	fmt.Println("usage:", usageInfo[subcmd])
}

func PrintUsages() {
	fmt.Println("Usage: octl <subcmd> <args...>\n[subcmd]:")
	for subcmd, usage := range usageInfo {
		fmt.Printf("- %s\n    usage: %s\n", subcmd, usage)
	}
}
