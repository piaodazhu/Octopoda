![Octopoda](assets/logo.gif)
# Octopoda

🐙 **Octopoda** is a lightweight multi-node scenario management platform. It's not a lightweight K8S. It is originally designed for managing Lab101's ICN application scenarios (Obviously it can do more than that), which require the execution of commands on the node at the lower level of the system, such as inserting a kernel driver module. **Note that it not safe enough to deploy Octopoda in unfamiliar network environment.**

Features of Octopoda:
1. Simple topology with NAT reverse path.
2. Out-of-box.
3. Robust & auto retry & auto reboot.
4. Nodes status monitoring.
5. Customized, automated scenario deployment.
6. Scenario/Application version control.
7. Scenario/Application durability.
8. Centralized file management and distribution.
9. Centralized scripts execution.
10. Log management.
11. Fast SSH login.
12. Golang/C/Python SDK.

# Table of Contents
- [Octopoda](#octopoda)
- [Table of Contents](#table-of-contents)
- [Concepts](#concepts)
  - [Topology](#topology)
  - [SAN Model](#san-model)
- [Quick Start](#quick-start)
- [Octl Command Manual](#octl-command-manual)
  - [A. Node Information](#a-node-information)
    - [GET](#get)
    - [PRUNE](#prune)
    - [STATUS](#status)
  - [B. Workgroup](#b-workgroup)
    - [CONCEPT](#concept)
    - [PATH COMMAND](#path-command)
    - [MEMBERS COMMAND](#members-command)
  - [C. Command Exection](#c-command-exection)
    - [RUN](#run)
    - [XRUN](#xrun)
  - [D. File Distribution](#d-file-distribution)
    - [UPLOAD](#upload)
    - [DOWNLOAD](#download)
  - [E. Fast SSH](#e-fast-ssh)
    - [SET](#set)
    - [LS](#ls)
    - [DEL](#del)
    - [LOGIN](#login)
  - [F. Scenario Deployment](#f-scenario-deployment)
    - [CREATE](#create)
    - [REPO](#repo)
    - [APPLY](#apply)
      - [a. About Target](#a-about-target)
      - [b. About Order](#b-about-order)
      - [c. About Path](#c-about-path)
      - [d. Script Environment Variables](#d-script-environment-variables)
    - [NODEAPP](#nodeapp)
  - [G. Version Control](#g-version-control)
    - [VERSION](#version)
    - [RESET](#reset)
  - [H. HttpsNameServer](#h-httpsnameserver)
    - [APIs](#apis)
    - [CertGen](#certgen)
  - [I. Online Upgrade - Pakma](#i-online-upgrade---pakma)
- [Scenario Example](#scenario-example)
- [Octl SDK](#octl-sdk)
  - [Golang SDK](#golang-sdk)
  - [C/C++ SDK](#cc-sdk)
  - [Python SDK](#python-sdk)


# Concepts

## Topology
```
                                            +-----------------+ 
                                            |    NameServer   | 
                                            +-----------------+   HTTP 
                                             | HTTPS         |-------------+  
                     SSH .---------->-----+--|--<-reverse----|-<-+------<--|--+------ ... 
                         |                |  |               |   |         |  |  
  +---------+  CLI  +--------+  HTTPS   +-x-------+   TLS   +----------+ +----------+  
  | TermUser| <===> |  Octl  | <--+---> |  Brain  | <-----> | Tentacle | | Tentacle | ...  
  +---------+       +  ----  +    |     +---------+         +----------+ +----------+  
  +---------+  SDK  | Go/C/Py| <--+          |HTTP               |             | 
  |Developer| <===> | Client |          +---------+         +---------+   +---------+ 
  +---------+       +--------+          |  Pakma  |         |  Pakma  |   |  Pakma  | 
                                        +---------+         +---------+   +---------+ 
                                       \-----------/       \-----------------------------/ 
                                         Master Node           Controlled Networks 
```
## SAN Model
**SAN** model is the working model of Octopoda: `S` stands for `Scenario`, A stands for `Application` and `N` stands for `Node`. The current model has the following features:
- An Octopoda network can manage multiple scenarios.
- Each scenario is made of multiple applications.
- Each application can be run on multiple nodes.
- Each application can be identified uniquely by (application name, scenario name)
- Each node can run multiple application.
- Each application running on a specific node is called a **NodeApp**, and it can be identified uniquely by (nodename, application name, scenario name).
- The granularity of version control are scenario and NodeApp. **Version updates of NodeApps automatically trigger version updates of the scenario they belong to.**

# Quick Start
```sh
# 1 Generate keys and certificates
cd httpNameServer
bash ./CertGen.sh ca helloCa # generate your ca key and cert ( keep your ca.key very safe! )
bash ./CertGen.sh server helloNs 1.1.1.1 # for your httpsNameServer
bash ./CertGen.sh client helloBrain # for your Brain
bash ./CertGen.sh client helloTentacle # for your Tentacle 

cp ca.pem server.key server.pem /etc/octopoda/cert/ # copy to your httpsNameServer
cp ca.pem client.key client.pem /etc/octopoda/cert/ # copy to your Brain and Tentacle

# 2 Install httpsNameServer on a machine
sudo systemctl start redis-server
tar -Jxvf httpns_v1.5.2_linux_amd64.tar.xz
cd httpns_v1.5.2_linux_amd64
# then modify httpns.yaml if you want 
# run it foreground
sudo ./httpns -p
# or install it and start background
sudo bash setup.sh

# 3 Install Brain on a machine
sudo systemctl start redis-server
tar -Jxvf brain_v1.5.2_linux_amd64.tar.xz
cd brain_v1.5.2_linux_amd64
# then modify brain.yaml if you want (httpsNameServer, name, nic is important)
# run it foreground
sudo ./brain -p
# or install it and start background
sudo bash setup.sh

## 4 Set root workgroup password on Brain machine. (For step 7)
redis-cli
127.0.0.1:6379> set info: yourpass

# 5 Install Tentacle on a machine
tar -Jxvf tentacle_v1.5.2_linux_amd64.tar.xz
cd tentacle_v1.5.2_linux_amd64
# then modify tentacle.yaml if you want (httpsNameServer, name, brain is important)
# run it foreground
sudo ./tentacle -p
# or install it and start background
sudo bash setup.sh

# 6 Install Pakma on your Brain or Tentacle machine (optional, only for online upgrade)
tar -Jxvf pakma_v1.5.2_linux_amd64.tar.xz
cd pakma_v1.5.2_linux_amd64
# make sure pakma is installed after Brain or Tentacle
# install it for your Brain
sudo bash setup.sh brain
# or install it for your Tentacle
sudo bash setup.sh tentacle

# 7 Install Octl
cd octl_v1.5.2_linux_amd64
# then modify octl.yaml. (workgroup.root="", workgroup.password="yourpass")
sudo cp octl.yaml /etc/octopoda/octl/
sudo cp octl /usr/local/bin/

# 8 Hello World
$ octl node get
# {
#   "nodes": [
#     {
#       "name": "pi0",
#       "addr": "192.168.1.4",
#       "state": "online",
#       "delay": "3ms",
#       "online_time": "2m42.483697064s"
#     }
#   ],
#   "total": 1,
#   "active": 1,
#   "offline": 0
# }

$ octl node get | grep name | awk '{print $2}' | sed 's/"//g' | sed -z 's/\n/ /g' | sed 's/,//g'
# you may get: pi02 pi05 pi06 pi08
```

# Octl Command Manual

## A. Node Information
### GET

> `usage: octl node get [-sf <statefilter>] [[ALL] | <node1> <@group1> ...]`

With this subcmd we can get some basic information of all nodes or given nodes. With optional flag `-sf`, you we can define a state filter such as `online`, `offline` to filter nodes.

<!-- 
- Basic information of all scenarios in the network or detailed information of a given scenario.
- Basic informations of all apps on the given node or detailed information of a given app on the given node. -->

### PRUNE
> `usage: octl node prune [ALL | <node1> <@group1> ...]`

With this subcmd we can prune dead nodes within given nodes.

### STATUS

> `usage: octl node status [[ALL] | <node1> <@group1> ...]`

With this subcmd we can get the running status of nodes or a given node, such as:
- CPU Load.
- Memory Used/Total.
- Disk Used/Total.
- Other Status.
<!-- 
### LOG

> `usage: octl log [<node>|brain] [-l<maxline>] [-d<maxday>]`

With this subcmd we can get the running log of brain or a given node. The argument `l` means max lines need to be read, and argument `d` means max days before today need to be read.

Default `l` is 30 and default `d` is 0, means latest 30 lines of logs will be return. -->

## B. Workgroup

### CONCEPT

Workgroup are an Octopoda mechanism that supports resource isolation at node granularity, hierarchical device authorization and referencing multiple node names with a group name. Workgroups are organized in a **tree structure**, with each node having a unique path, a non-empty set of node names, and a collection of subworkgroups. 

For example, workgroup with path `/room1/alice` and node names set `(pi1,pi2,pi3)` will never be aware of `pi4` even if it is in the same Octopoda network. And group `/room1` can list, add and remove members in/to/from `/room1/alice`. And `/room1/alice` can create `/room1/alice/g1` with node names set `(pi1,pi2)`. 

The **relative** group name can be used to reference its node names set. For example, if **currentPath=/room1**, then `octl cmd run 'uname -a' @alice pi4` is equivalent to `octl cmd run 'uname -a' pi1 pi2 pi3 pi4`.

Workgroup is a mechanism forshould be configured in `octl.yaml`.

```yaml
workgroup:
  root: "grouppath"
  password: "password"
  currentPathFile: "/etc/octopoda/octl/.curPath.yaml"
```

### PATH COMMAND
> `usage: octl wg pwd`
> 
> `usage: octl ls [<grouppath>]`
> 
> `usage: octl cd [<grouppath>]`
> 
> `usage: octl wg grant <grouppath> <password>`


### MEMBERS COMMAND
> `usage: octl wg get [<grouppath>]`
> 
> `usage: octl wg add <grouppath> [<node1>] [@<group1>] ...`
> 
> `usage: octl wg rm <grouppath> [[<node1>] [@<node1>] ...]`

## C. Command Exection

### RUN

> `usage: octl cmd run [-ta] [[-c] 'cmd' | -bg 'cmd' | -ss 'shellScript'] <node1> <@group1> ...`

With this subcmd we can run a command or a script on given nodes. For running a forground command, we can use flag `-c 'cmd'` or only `'cmd'`. As for blocking command, we need to run it background, so we can use flag `-bg 'cmd'`. For running a script, we need to specify the complete filepath of the script with flag `-ss 'shellScript'`.

The `-ta` (time align) flag is optional. When this enable `-ta` flag, all target nodes need to make sure that their process will start at the same moment as much as possible (with the help of Octopoda Simple Time Protocol). Since there is jitter in the delay, if the node cannot guarantee it, it will reject the execution.

### XRUN

> `usage: octl cmd xrun [-ta] [[-c] 'cmd' | -bg 'cmd' | -ss 'shellScript'] [-d <delayseconds>] <node1> <@group1> ...`

The difference between subcmd `run` and `xrun` is that `xrun` will not execute and return the result immediately, instead it just load the process and trigger it after a delay (specified with `-d` flag).

## D. File Distribution

### UPLOAD

> `usage: octl file upload [-f] <localFileOrDir> <targetDir>  [<node1> ...] [<@group1> ...]`

With this subcmd we can upload a file or a whole directory to the given nodes' storage. The `-f` flag means `targetDir` will be created by force if it not exists on target nodes. `targetDir` support these **path variable**:
- `@root`: the home directory of root user (/root/).
- `@workspace`: the root of workspace, configured by `tentacle.yaml` or `brain.yaml`.
- `@log`: log directory, configured by `tentacle.yaml` or `brain.yaml`.
- `@fstore`: file storage directory, configured by `tentacle.yaml` or `brain.yaml`.
- `@pakma`: pakma directory, configured by `tentacle.yaml` or `brain.yaml`.

### DOWNLOAD

> `usage: octl file download FileOrDir [localDir] <node1>`

With this subcmd we can download file or directory from under `FileOrDir` from brain or a given node to `localDir`. `FileOrDir` support these **path variable**:
- `@root`: the home directory of root user (/root/).
- `@workspace`: the root of workspace, configured by `tentacle.yaml` or `brain.yaml`.
- `@log`: log directory, configured by `tentacle.yaml` or `brain.yaml`.
- `@fstore`: file storage directory, configured by `tentacle.yaml` or `brain.yaml`.
- `@pakma`: pakma directory, configured by `tentacle.yaml` or `brain.yaml`.

## E. Fast SSH

### SET

> `usage: octl ssh set <node1>`

Set SSH login information binding with a name. 

### LS

> `usage: octl ls`

List all SSH services in this network.

### DEL

> `usage: octl ssh del <name>`

Delete SSH login information binding with a name. 

### LOGIN

> `usage: octl ssh login <name>`

Directly login the host binding with a name via SSH.

## F. Scenario Deployment
### CREATE

> `usage: octl scen create <scenario> [with <app1> <app2> ...]`

With this subcmd we can generate a template folder for a scenarios, which can be edited and then used for `apply` subcmd.

### REPO

> `usage: octl scen repo [clone|push] <scen> [-u <username>]`

With this subcmd we can clone a scenario deployment folder from a git service, or push modification to a remote git service.

Note that the git service URL and Auth information should be configured in `octl.yaml`. It is recommanded that deploying a self-hosted git service with `gogs`.

### APPLY

> `usage: octl scen apply <scenario> [target] -m "your message"`

With this subcmd we can create, delete, run a scenario. The information required for scenario deployment is defined in the `<scenario>` folder, with a deployment.yaml configuration file and application subfolders. Below is a typical deployment file. 

```yaml
# a simple example
name: helloWorld
description: "a scenario to print hello world"
applications:
-
  name: "helloPrinter"
  # scripts should be found under this path of current host
  scriptpath: "hello/scripts/"
  sourcepath: "hello/src/"
  description: "print hello to hello.txt"
  nodes:
    - 'pi-24'
  script:
    -
      # it will prepare a file hello.txt and write first line `-->prepared`
      target: prepare
      file: "prepare.sh"
      order: 3
    -
      target: start
      # it will write second line `hello` and print `hello`
      file: "run.sh"
      order: 1
    -
      # it will do nothing but print `stop`
      target: stop
      file: "stop.sh"
      order: 1
    -
      # it will delete the file hello.txt
      target: purge
      file: "purge.sh"
      order: 1
    -
      # it will delete the file hello.txt
      target: date
      file: "date.sh"
      order: 1
-
  name: "worldPrinter"
  scriptpath: "world/scripts/"
  sourcepath: "world/src/"
  description: "print world to world.txt"
  nodes:
    - 'pi-240'
  script:
    -
      # it will prepare a file world.txt and write first line `-->prepared`
      target: prepare
      file: "prepare.sh"
      order: 1
    -
      target: start
      # it will write second line `world` and print `world`
      file: "run.sh"
      order: 2
    -
      # it will do nothing but print `stop`
      target: stop
      file: "stop.sh"
      order: 2
    -
      # it will delete the file world.txt
      target: purge
      file: "purge.sh"
      order: 2
-
  name: "helloWorldPrinter"
  scriptpath: "helloworld/scripts/"
  sourcepath: "helloworld/src/"
  description: "print helloworld to helloworld.txt"
  nodes:
    - 'pi-24'
    - 'pi-240'
  script:
    -
      # it will prepare a file helloworld.txt and write first line `-->prepared`
      target: prepare
      file: "prepare.sh"
      order: 2
    -
      target: start
      # it will write second line `helloworld` and print `helloworld`
      file: "run.sh"
      order: 3
    -
      # it will do nothing but print `stop`
      target: stop
      file: "stop.sh"
      order: 3
    -
      # it will delete the file helloworld.txt
      target: purge
      file: "purge.sh"
      order: 3
    -
      # it will delete the file hello.txt
      target: date
      file: "date.sh"
      order: 3
```

Except for this deployment file, some shell script should be prepared. It will be executed at the right time. And source files are optional, it will be copied to the right nodes when the scenario is prepared, and can be used as material for script execution.

#### a. About Target

`target` means the target of a scenario operation. If subcmd `apply` run with `target`, scripts corresponding to this `target` will be executed. There are three types of target:
- **SPECIAL** target: `prepare`, `purge`, `commit`
- **NORMAL** target: `start`, `stop`
- **CUSTOMIZED** target: defined by user.

If subcmd `apply` run with **SPECIAL** target, Octopoda will not just run the corresponding scripts (`commit` will not execute any scripts). **SPECIAL** target and **MNORMAL** target must be implemented for each application. And **CUSTOMIZED** is defined by user.

When running subcmd `apply` in command line, target `default` or `(empty)` means `prepare + start`. For target `commit`, `-m {message}` is necessary. 

#### b. About Order

`order` is a optional field. It can be used when a scenario is executing a `target`, the scripts corresponding to that `target` need to be executed in a specific order. The smaller `order`, the scripts will be executed in order of order from smallest to largest. For the same `order`, the execution sequence is random. Default `order` is 0. 

#### c. About Path

`scriptpath`: The root path of the scripts. All mentioned scripts must be available in admin client's storage.
`sourcepath`: The root path of the sourcefiles. All files or directories under this path will be copied to the right nodes before any script is executed.

#### d. Script Environment Variables

Some environment variables are predefined when any script is executed. It can be directly referred in any script. Current Octopoda support these:
- `OCTOPODA_NODENAME`: the name of the node who executes this script.
- `OCTOPODA_CURRENTDIR`: work directory of the node who executes this script.
- `OCTOPODA_FILENAME`: file name of this script.
- `OCTOPODA_OUTPUT`: output file of this script.
- `OCTOPODA_APP`: current application name of this script.
- `OCTOPODA_SCENARIO`: current scenario name of this script.
- `OCTOPODA_ROOT`: system super user home directory of the node who executes this script.
- `OCTOPODA_WORKSPACE`: workspace root directory of the node who executes this script.
- `OCTOPODA_FSTORE`: storage directory of the node who executes this script.
- `OCTOPODA_LOG`: log directory of the node who executes this script.
- `OCTOPODA_PAKMA`: pakma directory of the node who executes this script.

Note that in scripts, output to stdout (for example, `echo "done"`) won't work. If some output information need to be collected and shown in execution results, we have to **append them to `OCTOPODA_OUTPUT`**. (for example, `echo "done" >> $OCTOPODA_OUTPUT`)

Customized environment variables are also supported. They can be defined in tentacle.yaml and can be directly referred in script or command.

### NODEAPP

> `usage: octl napp get <node> [[ALL] | <app>@<scen>]`

With this subcmd we can get basic informations of all apps on the given node or detailed information of a given app on the given node.

## G. Version Control
### VERSION

> `usage: octl scen version <scen> [-o <offset>] [-l <limit>]`
> `usage: octl napp version <node> <app>@<scen> [-o <offset>] [-l <limit>]`

With this subcmd we can get a version list of a given scenario, or a given app on a given node. Each version consists of version hash code, committed message, committed timestamp and other basic information.

Current Octopoda only support version list so `branch` is not supported. However, as a tool for scenario deployment, version list is enough in most case. Complex version control should occur primarily in the development phase.

`.gitignore` is supported.

Changes will be committed when applying target `prepare` and `commit`.

### RESET

> `usage: octl scen reset <scen> -v <version> -m "your message"`
> `usage: octl napp reset <node> <app>@<scen> -v <version> -m "your message"`

With this subcmd we can set a given scenario, or a given app on a given node to a given historical version. **If a scenario is reset, all relative apps on corresponding nodes will be reset. If an app on a given node is reset, the corresponding scenario will evolved into a new version. That's the rule**

The argument `version` need us to specify the prefix of the version hash code, whose length is at least 3 char.

Note that `reset` will not really let the version list back to a history version, but actually like `revert`. If we set `A->B->C` to `A`, the version list will become `A->B->C->A`, not `A`. And we won't lost version `C`.

Hot reset is not supported in current Octopoda. Therefore, stop the running scenario service before the reset, the start the running service after the reset.

<!-- ### FIX

> `usage: octl fix scenario <scen>`

With this subcmd we can manually fix the version file of a given scenario. When the actual version of the application in the scenario does not match the version in the version file, this subcmd may help. 

There is no need to run this subcmd in most cases, because fix will also be periodically executed by a goroutine. -->


## H. HttpsNameServer

### APIs

see `httpNameServer/api_doc.md`

### CertGen

Fast script to generate keys and certificates for CA, server and client.

## I. Online Upgrade - Pakma

> `usage: octl pakma [state|install <version>|upgrade <version>|confirm|cancel|downgrade|history|clean] [<brain>|<node1>|<node2>|...] [-t<timestr>] [-l<limit>]`

Pakma is a localhost http service to support tentacle/brain online upgrade, downgrade. So far, HttpsNameServer must be enabled to provide [release package](https://github.com/piaodazhu/Octopoda/releases).

Pakma's basic mechanism is **stable-(upgrade)-preview-(confirm/cancel/timeout)-stable**. If the upgraded version of tentacle/brain can not work well, pakma will automatically rollback it to the old stable version after 2 minites timeout.

- example1: octl pakma install 1.5.1 brain pi0 pi1 pi2   (first upgrade, use install)
- example2: octl pakma state brain pi0 pi1 pi2
- example3: octl pakma upgrade 1.5.1 brain pi0 pi1 pi2   (after first upgrade, use upgrade)

# Scenario Example
See an example in `./octl/example/helloWorld`. The file `deployment.yaml` defines scenario called `helloWorld`. This scenario consists of 3 application, running on 2 nodes, with some targets. You can manage `helloWorld` scenario with octl:

```sh
cd ./octl
# SPECAIL target `prepare` will copy all of the sourcepath files to corresponding nodes, 
# then run all scripts of target `prepare` on corresponding nodes. 
./octl scen apply example/helloWorld prepare -m "Prepare a new scenario"

# NORMAL target `start` will run all scripts of target `start` on corresponding nodes. 
./octl scen apply example/helloWorld start -m "Start run a scenario"

# NORMAL target `stop` will run all scripts of target `stop` on corresponding nodes. 
./octl scen apply example/helloWorld stop -m "Stop run a scenario"

# CUSTOMIZED target `date` will run all scripts of target `date` on corresponding nodes. 
./octl scen apply example/helloWorld date -m "Append current date to file"

# SPECIAL target `commit` will commit all changes on corresponding nodes. 
./octl scen apply example/helloWorld commit -m "Append current date to file"

# see details of this scenario
./octl scen get scenario helloWorld

# list all history versions of this scenario. 
./octl scen version scenario helloWorld

# reset this scenario to a history version. 
./octl scen reset scenario helloWorld -v d8ef -m "Reset to a history version"

# SPECAIL target `purge` will run all scripts of target `purge` on corresponding nodes,
# then remove all nodeapp directories on corresponding nodes.
./octl scen apply example/helloWorld purge
```

The nodeApps will be installed under `{tentacle.yaml:workspace.root}/app@scen` on corresponding nodes.

# Octl SDK

We provide Golang/C/Python octl SDK to support integrating Octopoda's capability into your code.
## Golang SDK
APIs are exported in `octl/sdk/sdk.go`. For using them in Golang program, you should import this module first:
```go
import octl "github.com/piaodazhu/Octopoda/octl/sdk"
```
or `go get`:
```bash
go get "github.com/piaodazhu/Octopoda/octl/sdk"
```
Before programming with those APIs, make sure `Init()` has been called first and only once.
```go
// Init only once
if err := octl.Init("/your/path/of/octl_conf.yaml"); err != nil {
  // something wrong
}

nodesInfo, err := octl.NodeInfo([]string{"node1", "node2", "node3"})
if err != nil {
  // something wrong
}

```
Some usage examples can be found in `octl/sdk/sdk_test.go`.

## C/C++ SDK
APIs are exported in `octl/sdk/coctl/coctl.h`. For using them in C/C++ program, you should include `wrapper.h` and `coctl.h`, and link `libcoctl.a` or `libcoctl.so`.

When programming with those APIs, you should call `octl_init` first. Those APIs take C style with input/output pointer, and **all input array of pointers must be freed manually (See examples)**.

Some usage examples can be found in `octl/sdk/coctl/test.c`.

## Python SDK
APIs are exported in `octl/sdk/pyoctl/pyoctl.py`. For using them in Python program, you should import `pyoctl`, and make sure you have `libcoctl.so` or `libcoctl.dll`.
```python
try:
  octl = pyoctl.OctlClient("/path/of/libcoctl.so", "/path/of/octl_conf.yaml")
  results = octl.run_command("uname -a", ['node1', 'node2'])
    for result in results:
      print(result)
except pyoctl.OctlException as e: # capture the exception of octl
  print(e)

```

Some usage examples can be found in `octl/sdk/pyoctl/test.py`.

 