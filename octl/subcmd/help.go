package subcmd

import (
	"fmt"
	"sort"
	"strings"
)

type UsageEntry struct {
	ClassSerial string
	Class       string
	Operation   string
	Usage       string
	Example     string
}

var UsageList []UsageEntry
var Usages = `

a:Node:get:octl node get [-sf <statefilter>] [[ALL] | <node1> <@group1> ...]:octl node get -sf online ALL
a:Node:prune:octl node prune [ALL | <node1> <@group1> ...]:octl node prune @mygroup
a:Node:status:octl node status [[ALL] | <node1> <@group1> ...]:octl node status pi0 pi1

b:Workgroup:pwd:octl wg pwd:octl wg pwd
b:Workgroup:cd:octl wg cd [<grouppath>]:octl wg cd mygroup
b:Workgroup:ls:octl wg ls [<grouppath>]:octl wg ls .
b:Workgroup:grant:octl wg grant <grouppath> <password>:octl wg grant subgroup 123456
b:Workgroup:get:octl wg get [<grouppath>]:octl wg get mygroup
b:Workgroup:add:octl wg add <grouppath> <node1> @<node1>:octl wg add g1 pi1 pi2
b:Workgroup:rm:octl wg rm <grouppath> <node1> @<node1>:octl wg rm g1 pi1

c:Command:run:octl cmd run [-ta] [[-c] 'cmd' | -bg 'cmd' | -ss 'shellScript'] <node1> <@group1> ...:octl cmd run -ta -c 'uname -a' @mygroup pi3
c:Command:xrun:octl cmd xrun [-ta] [[-c] 'cmd' | -bg 'cmd' | -ss 'shellScript'] [-d <delayseconds>] <node1> <@group1> ...:octl cmd xrun -ta -c 'uname -a' -d 10 pi0 pi2

d:File:upload:octl file upload [-f] <localFileOrDir> <targetDir>  [<node1> ...] [<@group1> ...]:octl upload './pictures' 'Pictures' pi0 pi1 pi2
d:File:download:octl file download FileOrDir [localDir] <node1>:octl file download 'Pictures' './' pi2

e:SSH:login:octl ssh login <node1>:octl ssh login pi0
e:SSH:set:octl ssh set <node1>:octl ssh set pi0
e:SSH:del:octl ssh del <node1>:octl ssh del pi0
e:SSH:ls:octl ssh ls:octl ssh ls

f:Scenario:create:octl scen create <scen> [with <app1> <app2> ...]:octl scen create ChatScen with alice bob
f:Scenario:repo:octl scen repo [clone|push] <scen> [-u <username>]:octl scen repo clone ChatScen -u mike
f:Scenario:apply:octl scen apply <scen> [target] -m "your message":octl scen apply ChatScen prepare -m "prepare my scenario"
f:Scenario:version:octl scen version <scen> [-o <offset>] [-l <limit>]:octl scen version scenario ChatScen -o 0 -l 10
f:Scenario:reset:octl scen reset <scen> -v <version> -m "your message":octl scen reset scenario ChatScen -v b698 -m "back to yesterday"

g:NodeApp:get:octl napp get <node> [[ALL] | <app>@<scen> [-o <offset>] [-l <limit>]]:octl napp get pi0 alice@ChatScen
g:NodeApp:reset:octl napp reset <node> <app>@<scen> -v <version> -m "your message":octl napp reset pi0 alice@ChatScen -v b698 -m "back to yesterday"

h:PAcKage MAnager:pakma:octl pakma [state|install <version>|upgrade <version>|confirm|cancel|downgrade|history|clean] [<brain>|<node1>|<@group1>|...] [-t<timestr>] [-l<limit>]:octl pakma upgrade 1.5.1 brain @upgradable
`

func InitUsage() {
	UsageList = []UsageEntry{}
	lines := strings.Split(Usages, "\n")
	for _, l := range lines {
		if len(l) == 0 {
			continue
		}
		fields := strings.Split(l, ":")
		if len(fields) != 5 {
			panic("InitUsage Error Because Usage Format is Bad")
		}
		UsageList = append(UsageList, UsageEntry{
			ClassSerial: fields[0],
			Class:       fields[1],
			Operation:   fields[2],
			Usage:       fields[3],
			Example:     fields[4],
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
	if len(args) == 1 {
		PrintUsage(args[0], "")
		return
	}

	usageClassMap := map[string][]UsageEntry{}
	for _, usageEntry := range UsageList {
		usageClassMap[usageEntry.Class] = append(usageClassMap[usageEntry.Class], usageEntry)
	}

	usageClassList := [][]UsageEntry{}
	for _, value := range usageClassMap {
		usageClassList = append(usageClassList, value)
	}

	sort.Slice(usageClassList, func(i, j int) bool { return usageClassList[i][0].ClassSerial < usageClassList[j][0].ClassSerial })

	fmt.Println("Usage: octl <subcmd> <args...>\n[subcmd]:")
	for _, subList := range usageClassList {
		fmt.Printf("%s. %s\n", subList[0].ClassSerial, subList[0].Class)
		for _, entry := range subList {
			fmt.Printf("    - %s\n        usage: %s\n        example: %s\n",
				entry.Operation, entry.Usage, entry.Example)
		}
		fmt.Println()
	}
}
