## deploy scienario
```sh
octl apply <deployfile> <target:prepare,start,stop,purge>
```
Example of deployfile
```yaml
name: scienario1
desciption: "hello world"
scriptpath: "./script/"
applications:
- 
	name: "icn_provider"
	nodes:
		- pi1
		- pi2
	prepare:
		script: "prepare.sh"
	start:
		script: "run.sh"
	stop:
		script: "stop.sh"
	purge:
		script: "purge.sh"
- 
	name: "icn_consumer"
	nodes:
		- pi3
		- pi4
	prepare:
		script: "prepare2.sh"
	start:
		script: "run2.sh"
	stop:
		script: "stop2.sh"
	purge:
		script: "purge2.sh"
```

## get infomation
```sh
octl get nodes
octl status nodes
octl get node <name>
octl get scienarios
octl get scienario <name>
```
Example of output
```sh
$ octl get nodes
name	platform	status	ip	port
pi2		armv7	active  192.168.3.240	1234 

$ octl status nodes
name 	shortload	longload	memused	memtatal	diskused	disktotal	verison	age
pi2	2.4  0.9	1234  2333	4096	11111	65536	3	40m

$ octl get node pi2
name: p2
shortload: 2.4
longload: 0.9
memused: 2333
memtotal: 4096
diskused: 11111
disktotal: 65536
version: 3
age: 40m
applications:
	app1:
		version: 3
		status: running
	app2:
		version: 1
		status: stop

```

## get log
```sh
octl log node <name>
octl log scienarios
octl log scienario <name>
```

## get history versions
```sh
octl version nodes
octl version node <name>
octl version scienarios
octl version scienario <name>
```

## switch the versions
```sh
octl reset nodes <name1,name2,...> <version> message <message>
octl reset node <name> <version> message <message>
octl reset scienarios <name> <version> message <message>
```

## enter shell
```sh
octl shell node <name>
```

## send files
```sh
octl sendfile <filelist.txt> <nodelist.txt>
```

## purge dead nodes
```sh
octl purge nodes
```
