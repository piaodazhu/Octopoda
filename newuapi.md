#1 Node				2 (+2)
	[CHK]	- node get
			[TODO:	Version 1]
			usage: octl node get [[ALL] | <node1> <@group1> ...]
			[TODO:	Version 2]
			usage: octl node get [-on|-dis|-off] [[ALL] | <node1> <@group1> ...]
	[CHK]	- node prune
			[TODO:	Version 1]
			usage: octl node prune [ALL | <node1> <node2> ...]
	[HIDE]	- node log
	[HIDE]	- node status

#2 Group			3
	[CHK]	- group get
			[TODO:	Version 1]
			usage: octl group get [[ALL] | <group>]
			[Later:	Version 2]
			usage: octl group get [[ALL] | <group> ...]
	[CHK]	- group set
			[TODO:	Version 1]
			usage: octl group set [-n] <group> <nodes...>
			# -n == --no-check
	[CHK]	- group del
			[TODO:	Version 1]
			usage: octl group del [<group>]
			[Later:	Version 2]
			usage: octl group del [<group> ...]

#3 CMD				2
	[CHK]	- cmd run
			[TODO:	Version 1]
			usage: octl run [-ta] [[-c] 'cmd' | -bg 'cmd' | -s 'shellScript'] <node1> <@group1> ...
			# -ta == --try-align
			# -cc == --common-command	(default)
			# -bg == --background
			# -ss == --shellscript
	[CHK]	- cmd xrun
			[TODO:	Version 1]
			usage: octl xrun [-ta] [[-c] 'command' | -bg 'command' | -s 'shellScript'] [-d <delayseconds>] <node1> <@group1> ...
#4 File				2
	[CHK]	- file upload
			[TODO:	Version 1]
			usage: octl upload [-f] <localFileOrDir> <targetDir> <node1> <node2> ...
			example: octl upload './pictures' '~/Pictures' pi0 pi1 pi2
			[Later:	Version 2]
			usage: octl upload [-f] <localFileOrDir> <targetDir> <node1> ... [<targetDir>] <nodeN> ...
	[CHK]	- file download
			[TODO:	Version 1]
			usage: octl download FileOrDir [localDir]
	
#5 SSH				4
	[CHK]	- ssh login
	[CHK]	- ssh set
	[CHK]	- ssh del
	[CHK]	- ssh ls

#6 Scenario			3 (+4)
	[??]	scen??

	[]	- scen apply
	[]	- scen create
	[CHK]	- scen get
			[TODO:	Version 1]
			usage: octl scen get [[ALL] | <scen>]

	[HIDE]	- scen repo
	[HIDE]	- scen version
	[HIDE]	- scen reset
	[HIDE]	- scen fix

#7 NodeApp			1
	[CHK]	- napp get
			[TODO:	Version 1]
			usage: octl napp get <node> [[ALL] | <app>@<scen>]

#8 PAcKage MAnager		8
	[]	- pakma state
	[]	- pakma install
	[]	- pakma upgrade
	[]	- pakma confirm
	[]	- pakma cancel
	[]	- pakma downgrade
	[]	- pakma history
	[]	- pakma clean