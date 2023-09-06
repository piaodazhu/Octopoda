package subcmd

import (
	"fmt"
	"strings"
)

type UsageEntry struct {
	Class   string
	Command string
	Usage   string
	Example string
}

var UsageList []UsageEntry
var Usages = `

Node:get:octl get [nodes|node <node>|scenarios|scenario <scen>|nodeapps <node>|nodeapp <node> <app>@<scen>]:octl get node pi0
Node:status:octl status [nodes|<node>|master]:octl status nodes
Node:prune:octl prune:octl prune
Node:log:octl log [master|<node>] [-l<maxline>] [-d<maxday>]:octl log pi0 -l50 -d2
Node:group:octl group [get <group>|del <group>|set <group> <nodes...>|set-nocheck <group> <nodes...>]:group set mygroup1 pi0 pi1 pi2

Scenario:create:octl create <scen> [with <app1> <app2> ...]:octl create ChatScen with alice bob
Scenario:apply:octl apply <scen> [target] -m "your message":octl apply ChatScen prepare -m "prepare my scenario"
Scenario:version:octl version [scenario <scen>|nodeapp <node> <app>@<scen>]:octl version scenario ChatScen
Scenario:reset:octl reset [scenario <scen>|nodeapp <node> <app>@<scen>]  -v <version> -m "your message":octl reset scenario ChatScen -v b698 -m "back to yesterday"
Scenario:fix:octl fix scenario <scen>:octl fix scenario ChatScen

File:upload:octl upload <localFileOrDir> <targetDir>:octl upload './pictures' './collections/img'
File:spread:octl spread <masterFileOrDir> <targetDir> <node1> <node2> ...:octl spread './collections/img' '~/Pictures' pi0 pi1 pi2
File:distrib:octl distrib <localFileOrDir> <targetDir> <node1> <node2> ...:octl distrib './pictures' '~/Pictures' pi0 pi1 pi2
File:tree:octl tree [store [master|<node>]|log [<node>|master]|nodeapp <node> app>@<scen>] [SubDir]:octl tree store pi0 '~/Pictures'
File:pull:octl pull [store [master|<node>]|log [<node>|master]|nodeapp <node> app>@<scen>] FileOrDir [localDir]:octl pull nodeapp pi0 alice@ChatScen './data/123.dat' './data'

SSH:ssh:octl ssh <anyname>:octl ssh pi0
SSH:setssh:octl setssh <anyname>:octl setssh pi0
SSH:getssh:octl getssh:octl getssh
SSH:delssh:octl delssh <anyname>:octl delssh pi0

Command:run:octl run [ '{<command>}' | '(<bgcommand>)' | <script> ] <node1> node2> ...:octl run '{ls ~/}' pi0

Upgrade:pakma:octl pakma [state|install <version>|upgrade <version>|confirm|cancel|downgrade|history|clean] [<master>|<node1>|<node2>|...] [-t<timestr>] [-l<limit>]:octl pakma upgrade 1.5.1 master pi0 pi1 pi2
`

func InitUsage() {
	UsageList = []UsageEntry{}
	lines := strings.Split(Usages, "\n")
	for _, l := range lines {
		if len(l) == 0 {
			continue
		}
		fields := strings.Split(l, ":")
		if len(fields) != 4 {
			panic("InitUsage Error Because Usage Format is Bad")
		}
		UsageList = append(UsageList, UsageEntry{
			Class:   fields[0],
			Command: fields[1],
			Usage:   fields[2],
			Example: fields[3],
		})
	}
}

func PrintUsage(subcmd string) {
	class := ""
	for i := range UsageList {
		if UsageList[i].Command == subcmd {
			fmt.Printf("- %s  class=%s\n    usage: %s\n    example: %s\n", 
				UsageList[i].Command, UsageList[i].Class, UsageList[i].Usage, UsageList[i].Example)
			class = UsageList[i].Class
		}
	}
	seealso := []string{}
	for i := range UsageList {
		if UsageList[i].Class == class {
			seealso = append(seealso, UsageList[i].Command)
		}
	}
	fmt.Printf("SEE ALSO: %s\n", strings.Join(seealso, ", "))

}

func PrintUsages(args []string) {
	if len(args) == 0 || len(args) > 1 {
		fmt.Println("Usage: octl <subcmd> <args...>\n[subcmd]:")
		for i := range UsageList {
			fmt.Printf("- %s  class=%s\n    usage: %s\n    example: %s\n", 
				UsageList[i].Command, UsageList[i].Class, UsageList[i].Usage, UsageList[i].Example)
		}
		return 
	}
	PrintUsage(args[0])
}
