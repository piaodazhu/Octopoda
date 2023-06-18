![Octopoda](assets/logo.gif)
# Octopoda

🐙 **Octopoda** is a lightweight multi-nodes scenario management platform. It's not a lightweight K8S. It is originally designed for managing Lab101's ICN application scenarios (Obviously it can do more than that), which require the execution of commands on the node at the lower level of the system, such as inserting a kernel driver module. **Note that it not safe enough to deploy Octopoda in unfamiliar network environment.**

Features of Octopoda:
1. Simple topology.
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

# Concepts

## Topology
```
                     SSH +--------------------------->-----------+------>-----+------ ...
     (Web)               |                                       |            | 
  +---------+  Run  +--------+   HTTP   +---------+   TCP   +----------+ +----------+ 
  |  Admin  | <===> |  Octl  | <------> |  Brain  | <-----> | Tentacle | | Tentacle | ... 
  +---------+       +--------+          +---------+         +----------+ +----------+  
 \----------------------------/        \-----------/       \-----------------------------/
          Admin Client                  Master Node           Controlled Networks
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

# Build & Try
```sh
# 1 build Tentacle
cd tentacle
go mod tidy
(GOOS=? GOARCH=? CGO_ENABLED=?) go build -o tentacle .

# 2 build Brain
cd brain
go mod tidy
(GOOS=? GOARCH=? CGO_ENABLED=?) go build -o brain .

# 3 build Octl
cd octl
go mod tidy
(GOOS=? GOARCH=? CGO_ENABLED=?) go build -o octl .

# 4 copy necessary files to corresponding nodes

# 5 edit corresponding configuration

# 6.1 run them in terminal (with -p to print log to Stdout)
./tentacle -p   # Should be root user or use sudo
./brain -p      # Should be root user or use sudo

# 6.2 you can also run them as deamon (Highly recommended)
bash setup.sh          # Install and run. Should be root user or use sudo
bash uninstall.sh      # Stop and Uninstall. Should be root user or use sudo

# 7 Manage the Octopoda network with Octl
./octl help     # get some help
./octl <subcmd> <args>
```

# Octl Command Manual

## A. Network Information
### GET

> `usage: octl get [nodes|node <node>|scenarios|scenario <scen>|nodeapps <node>|nodeapp <node> <app>@<scen>]`

With this subcmd we can get some basic information about current octopoda network, such as:
- Basic information of all nodes or detailed information of a given node.
- Basic information of all scenarios in the network or detailed information of a given scenario.
- Basic informations of all apps on the given node or detailed information of a given app on the given node.

### STATUS

> `usage: octl status [nodes|node <node>]`

With this subcmd we can get the running status of nodes or a given node, such as:
- CPU Load.
- Memory Used/Total.
- Disk Used/Total.
- Other Status.

### LOG

> `usage: octl log [master|node <node>] [l<maxline>] [d<maxday>]`

With this subcmd we can get the running log of master or a given node. The argument `l` means max lines need to be read, and argument `d` means max days before today need to be read.

Default `l` is 30 and default `d` is 0, means latest 30 lines of logs will be return.


## B. Scenario Deployment
### APPLY

> `usage: octl apply xx.yaml [target] -m "your message"`

With this subcmd we can create, delete, run a scenario. The information required for scenario deployment is defined in a yaml file. Below is a typical deployment file. 

```yaml
# a simple example
name: helloWorld
description: "a scenario to print hello world"
applications:
-
  name: "helloPrinter"
  # scripts should be found under this path of current host
  scriptpath: "./example/helloWorld/hello/scripts/"
  sourcepath: "./example/helloWorld/hello/src/"
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
  scriptpath: "./example/helloWorld/world/scripts/"
  sourcepath: "./example/helloWorld/world/src/"
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
  scriptpath: "./example/helloWorld/helloworld/scripts/"
  sourcepath: "./example/helloWorld/helloworld/src/"
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
- **SPECIAL** target: `prepare`, `purge`
- **NORMAL** target: `start`, `stop`
- **CUSTOMIZED** target: defined by user.

If subcmd `apply` run with **SPECIAL** target, Octopoda will not only run the corresponding scripts but also do something else. **SPECIAL** target and **MNORMAL** target must be implemented for each application. And **CUSTOMIZED** is defined by user.

When running subcmd `apply` in command line, target `default` or `(empty)` means `prepare + start`. For target `purge`, `-m {message}` is no necessary.

#### b. About Order

`order` is a optional field. It can be used when a scenario is executing a `target`, the scripts corresponding to that `target` need to be executed in a specific order. The smaller `order`, the scripts will be executed in order of order from smallest to largest. For the same `order`, the execution sequence is random. Default `order` is 0. 

#### c. About Path

`scriptpath`: The root path of the scripts. All mentioned scripts must be available in admin client's storage.
`sourcepath`: The root path of the sourcefiles. All files or directories under this path will be copied to the right nodes before any script is executed.

#### d. Script Environment Variables

Some environment variables are predefined when any script is executed. It can be directly referred in any script. Current Octopoda support these:
- `OCTOPODA_NODENAME`: the name of the node who executes this script.
- `OCTOPODA_CURRENTDIR`: work directory of the node who executes this script.
- `OCTOPODA_STOREDIR`: storage directory of the node who executes this script.
- `OCTOPODA_FILENAME`: file name of this script.
- `OCTOPODA_OUTPUT`: output file of this script.
- `OCTOPODA_APP`: current application name of this script.
- `OCTOPODA_SCENARIO`: current scenario name of this script.

Note that in scripts, output to stdout (for example, `echo "done"`) won't work. If some output information need to be collected and shown in execution results, we have to **append them to `OCTOPODA_OUTPUT`**. (for example, `echo "done" >> $OCTOPODA_OUTPUT`)

Customized environment variables are also supported. They can be defined in tentacle.yaml and can be directly referred in script or command.

## C. Version Control
### VERSION

> `usage: octl version [scenario <scen>|nodeapp <node> <app>@<scen>]`

With this subcmd we can get a version list of a given scenario, or a given app on a given node. Each version consists of version hash code, committed message, committed timestamp and other basic information.

Current Octopoda only support version list so `branch` is not supported. However, as a tool for scenario deployment, version list is enough in most case. Complex version control should occur primarily in the development phase.

`.gitignore` is supported.

### RESET

> `usage: octl reset [scenario <scen>|nodeapp <node> <app>@<scen>]  -v <version> -m <message>`

With this subcmd we can set a given scenario, or a given app on a given node to a given historical version. **If a scenario is reset, all relative apps on corresponding nodes will be reset. If an app on a given node is reset, the corresponding scenario will evolved into a new version. That's the rule**

The argument `version` need us to specify the prefix of the version hash code, whose length is at least 3 char.

Note that `reset` will not really let the version list back to a history version, but actually like `revert`. If we set `A->B->C` to `A`, the version list will become `A->B->C->A`, not `A`. And we won't lost version `C`.

Hot reset is not supported in current Octopoda. Therefore, stop the running scenario service before the reset, the start the running service after the reset.

### FIX

> `usage: octl fix scenario <scen>`

With this subcmd we can manually fix the version file of a given scenario. When the actual version of the application in the scenario does not match the version in the version file, this subcmd may help. 

There is no need to run this subcmd in most cases, because fix will also be periodically executed by a goroutine.

## D. File Distribution

### UPLOAD

> `usage: octl upload <localFileOrDir> <targetDir>`

With this subcmd we can upload a file or a whole directory to the storage of master. The argument `targetDir` is actually a relative path, so don't hesitate to set `targetDir` to the root directory `/`.

### SPREAD

> `usage: octl spread <masterFileOrDir> <targetDir> <node1> <node2> ...`

With this subcmd we can spread a file or a whole directory from the master's storage to the given nodes' storage. The argument `targetDir` is actually a relative path, so don't hesitate to set `targetDir` to the root directory `/`.


### DISTRIB

> `usage: octl distrib <localFileOrDir> <targetDir> <node1> <node2> ...`

With this subcmd we can distribute a file or a whole directory to the given nodes' storage. The argument `targetDir` is actually a relative path, so don't hesitate to set `targetDir` to the root directory `/`.

`distribute` can be considered as `upload` + `spread`.

### TREE

> `octl tree [store [master|<node>]|log [<node>|master]|nodeapp <node> <app>@<scen>] [SubDir]`

With this subcmd we can get all files information under `SubDir` on master or a given node. The pathtype `store`, `log` and `nodeapp` corresponds to the path configuration in the node configuration file. Current Octopoda not support print this files like a tree, it list all files instead. 

### PULL

> `octl pull [store [master|<node>]|log [<node>|master]|nodeapp <node> <app>@<scen>] FileOrDir [localDir]`

With this subcmd we can pull file or directory from under `SubDir` from master or a given node. The pathtype `store`, `log` and `nodeapp` corresponds to the path configuration in the node configuration file. 

## E. Script Execution

### RUN

> `usage: octl run [ '{<command>}' | '(<bgcommand>)' | <script> ] <node1> <node2> ...`

With this subcmd we can run a command or a script on given nodes. For running a command, we need to enclose the command in `'{}'`. As for blocking command, we need to run it in background, so we can enclose the command in `'()'`. For running a script, we need to specify the complete filepath of the script.

### SHELL

> `usage: octl shell <node>`

With this subcmd we can quickly login a given node via SSH. For some complex operation, this subcmd brings convenience. **Note that admin client should be able to connect to the node if we need to run this subcmd.**

However, this function in current Octopoda is not safe enough. If the network environment can't be trusted, don't fill the configuration file with actual usernamed and password.


# Scenario Example
See an example in `./octl/example/helloWorld`. The file `deployment.yaml` defines scenario called `helloWorld`. This scenario consists of 3 application, running on 2 nodes, with some targets. You can manage `helloWorld` scenario with octl:

```sh
cd ./octl
# SPECAIL target `prepare` will copy all of the sourcepath files to corresponding nodes, 
# then run all scripts of target `prepare` on corresponding nodes. 
./octl apply example/helloWorld/deployment.yaml prepare -m "Prepare a new scenario"

# NORMAL target `start` will run all scripts of target `start` on corresponding nodes. 
./octl apply example/helloWorld/deployment.yaml start -m "Start run a scenario"

# NORMAL target `stop` will run all scripts of target `stop` on corresponding nodes. 
./octl apply example/helloWorld/deployment.yaml stop -m "Stop run a scenario"

# CUSTOMIZED target `date` will run all scripts of target `date` on corresponding nodes. 
./octl apply example/helloWorld/deployment.yaml date -m "Append current date to file"

# see details of this scenario
./octl get scenario helloWorld

# list all history versions of this scenario. 
./octl version scenario helloWorld

# reset this scenario to a history version. 
./octl reset scenario helloWorld -v d8ef -m "Reset to a history version"

# SPECAIL target `purge` will run all scripts of target `purge` on corresponding nodes,
# then remove all nodeapp directories on corresponding nodes.
./octl apply example/helloWorld/deployment.yaml purge
```

