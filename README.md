![Octopoda](img/logo.gif)
# Octopoda

🐙 **Octopoda** is a lightweight multi-nodes scenario management platform. It's not a lightweight K8S. It is originally designed for managing Lab101's ICN application scenarios (Obviously it can do more than that), which require the execution of commands on the node at the lower level of the system, such as inserting a kernel driver module. **Note that it not safe enough to deploy Octopoda in unfamiliar network environment.**

Features of Octopoda:
1. Simple topology.
2. Robust & auto retry & auto reboot.
3. Nodes status monitoring.
4. Customized, automated scenario deployment.
5. Scenario/Application version control.
6. Scenario/Application durability.
7. Centralized file management and distribution.
8. Centralized scripts execution.
9. Log management.
10. Fast SSH login.

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
(GOOS=? GOARCH=?) go build -o tentacle .

# 2 build Brain
cd brain
go mod tidy
(GOOS=? GOARCH=?) go build -o brain .

# 3 build Octl
cd octl
go mod tidy
(GOOS=? GOARCH=?) go build -o octl .

# 4 copy necessary files to corresponding nodes

# 5 edit corresponding configuration

# 6.1 run them in terminal (with -p to print log to Stdout)
./tentacle -p
./brain -p

# 6.2 run them as deamon
bash setup.sh          # Install and run. Should be root user
bash uninstall.sh      # Stop and Uninstall. Should be root user

# 7 Manage the Octopoda network with Octl
./octl <subcmd> <args>
./octl <subcmd> <args> | jq   # jq makes the output JSON pretty
```

# Octl Command Manual

## octl get
- `get nodes`
- `get node {node name}`
- `get scenarios`
- `get scenario {scenario name}`
- `get nodeapps {node name}`
- `get nodeapp {node name} {app name}@{scenario name}`

## octl status
TBD

## octl apply
`octl apply {file name} {target} [-m {message}]`
- SPECIAL target: `prepare`, `purge`, `default`(=prepare+start), `(empty)`(=default)
- NORMAL target: `start`, `stop`
- CUSTOMIZED target: defined by user.
- For target `purge`, `-m {message}` is no necessary.

## octl version 
TBD

## octl reset
TBD

## octl log
TBD

## octl shell
TBD

## octl run
TBD

## octl filetree
TBD

## octl distrib
TBD

## octl fix
TBD

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
./octl get scenario helloWorld | jq

# list all history versions of this scenario. 
./octl version scenario helloWorld | jq

# reset this scenario to a history version. 
./octl reset scenario helloWorld -v d8ef -m "Reset to a history version" | jq

# SPECAIL target `purge` will run all scripts of target `purge` on corresponding nodes,
# then remove all nodeapp directories on corresponding nodes.
./octl apply example/helloWorld/deployment.yaml purge
```

