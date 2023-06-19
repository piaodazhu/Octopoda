package main

import (
	"flag"
	"fmt"
	"octl/config"
	"octl/nameclient"
	"octl/subcmd"
	"os"
)

var (
	BuildVersion string
	BuildTime    string
	BuildName    string
	CommitID     string
)

func main() {
	args := os.Args
	var conf string
	var askver bool
	var usage bool
	flag.BoolVar(&askver, "version", false, "tell version number")
	flag.StringVar(&conf, "c", "", "specify a configuration file")
	flag.BoolVar(&usage, "usage", false, "print subcommand usage")
	flag.Parse()

	if len(args) == 1 {
		fmt.Println("Octopoda Controlling Tool. Use '-usage', '-version', '-c'...")
		return
	}

	if askver {
		fmt.Printf("Octopoda Octl\nbuild name:\t%s\nbuild ver:\t%s\nbuild time:\t%s\nCommit ID:\t%s\n", BuildName, BuildVersion, BuildTime, CommitID)
		return
	}
	if usage {
		subcmd.PrintUsages()
		return 
	}
	if conf != "" {
		args = args[2:]
	}

	config.InitConfig(conf)
	nameclient.InitClient()
	switch args[1] {
	case "apply":
		subcmd.Apply(args[2:])
	case "get":
		subcmd.Get(args[2:])
	case "status":
		subcmd.Status(args[2:])
	case "log":
		subcmd.Log(args[2:])
	case "fix":
		subcmd.Fix(args[2:])
	case "version":
		subcmd.Version(args[2:])
	case "reset":
		subcmd.Reset(args[2:])
	case "shell":
		subcmd.Shell(args[2:])
	case "upload":
		subcmd.Upload(args[2:])
	case "spread":
		subcmd.Spread(args[2:])
	case "distrib":
		subcmd.Distrib(args[2:])
	case "tree":
		subcmd.FileTree(args[2:])
	case "pull":
		subcmd.Pull(args[2:])
	case "prune":
		subcmd.Prune(args[2:])
	case "run":
		subcmd.Run(args[2:])
	case "pakma":
		subcmd.Pakma(args[2:])
	case "help":
		subcmd.PrintUsages()
	default:
		fmt.Println("sub command not support")
	}
}
