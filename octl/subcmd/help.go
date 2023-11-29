package subcmd

import (
	"fmt"
	"strings"
)

type UsageEntry struct {
	Class   string
	Operation string
	Usage   string
	Example string
}

var UsageList []UsageEntry
var Usages = `

Node:get:octl node get [-sf <statefilter>] [[ALL] | <node1> <@group1> ...]:octl node get -sf online ALL
Node:prune:octl node prune [ALL | <node1> <@group1> ...]:octl node prune @mygroup
Node:status:octl node status [[ALL] | <node1> <@group1> ...]:octl node status pi0 pi1

Group:get:octl group get [[ALL] | <group>]:octl group get mygroup
Group:set:octl group set <group> [-n] [<node1> ...] [<@group1> ...]:octl group set mygroup pi0 pi2
Group:del:octl group del <group>:octl group del mygroup

Command:run:octl cmd run [-ta] [[-c] 'cmd' | -bg 'cmd' | -ss 'shellScript'] <node1> <@group1> ...:octl cmd run -ta -c 'uname -a' @mygroup pi3
Command:xrun:octl cmd xrun [-ta] [[-c] 'cmd' | -bg 'cmd' | -ss 'shellScript'] [-d <delayseconds>] <node1> <@group1> ...:octl cmd xrun -ta -c 'uname -a' -d 10 pi0 pi2

File:upload:octl file upload [-f] <localFileOrDir> <targetDir>  [<node1> ...] [<@group1> ...]:octl upload './pictures' 'Pictures' pi0 pi1 pi2
File:download:octl file download FileOrDir [localDir] <node1>:octl file download 'Pictures' './' pi2

SSH:login:octl ssh login <node1>:octl ssh login pi0
SSH:set:octl ssh set <node1>:octl ssh set pi0
SSH:del:octl ssh del <node1>:octl ssh del pi0
SSH:ls:octl ssh ls:octl ssh ls

Scenario:create:octl scen create <scen> [with <app1> <app2> ...]:octl scen create ChatScen with alice bob
Scenario:repo:octl scen repo [clone|push] <scen> [-u <username>]:octl scen repo clone ChatScen -u mike
Scenario:apply:octl scen apply <scen> [target] -m "your message":octl scen apply ChatScen prepare -m "prepare my scenario"
Scenario:version:octl scen version <scen>:octl scen version scenario ChatScen
Scenario:reset:octl scen reset <scen> -v <version> -m "your message":octl scen reset scenario ChatScen -v b698 -m "back to yesterday"

NodeApp:get:octl napp get <node> [[ALL] | <app>@<scen>]:octl napp get pi0 alice@ChatScen
NodeApp:reset:octl napp reset <node> <app>@<scen> -v <version> -m "your message":octl napp reset pi0 alice@ChatScen -v b698 -m "back to yesterday"

PAcKage MAnager:pakma:octl pakma [state|install <version>|upgrade <version>|confirm|cancel|downgrade|history|clean] [<brain>|<node1>|<group1>|...] [-t<timestr>] [-l<limit>]:octl pakma upgrade 1.5.1 brain pi0 pi1 pi2
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
			Operation: fields[1],
			Usage:   fields[2],
			Example: fields[3],
		})
	}
}

func PrintUsage(class, operation string) {
	for i := range UsageList {
		if UsageList[i].Class == class && UsageList[i].Operation == operation {
			fmt.Printf("- %s  class=%s\n    usage: %s\n    example: %s\n",
				UsageList[i].Operation, UsageList[i].Class, UsageList[i].Usage, UsageList[i].Example)
		}
	}
	seealso := []string{}
	for i := range UsageList {
		if UsageList[i].Class == class {
			seealso = append(seealso, UsageList[i].Operation)
		}
	}
	fmt.Printf("All operations of %s: %s\n", class, strings.Join(seealso, ", "))

}

func PrintUsages(args []string) {
	if len(args) == 0 || len(args) > 1 {
		fmt.Println("Usage: octl <subcmd> <args...>\n[subcmd]:")
		for i := range UsageList {
			fmt.Printf("- %s  class=%s\n    usage: %s\n    example: %s\n",
				UsageList[i].Operation, UsageList[i].Class, UsageList[i].Usage, UsageList[i].Example)
		}
		return
	}
	PrintUsage(args[0], "")
}
