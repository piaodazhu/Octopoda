package main

import (
	"flag"
	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/nameclient"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/octl/subcmd"
	"os"
)

// Build Information
var (
	BuildVersion string = "dev"
	BuildTime    string
	BuildName    string = "brain"
	CommitID     string = "snapshot"
)

func main() {
	output.EnableColor()
	output.EnableSpinner()
	
	subcmd.InitUsage()
	args := os.Args
	var conf string
	var askver bool
	var usage bool
	flag.BoolVar(&askver, "version", false, "tell version info")
	flag.StringVar(&conf, "c", "", "specify a configuration file")
	flag.BoolVar(&usage, "usage", false, "print subcommand usage")
	flag.Parse()

	if len(args) == 1 {
		output.PrintInfoln("Octopoda Controlling Tool. Use '-usage', '-version', '-c'...")
		return
	}

	if askver {
		output.PrintInfof("Octopoda Octl\nbuild name:\t%s\nbuild ver:\t%s\nbuild time:\t%s\nCommit ID:\t%s\n", BuildName, BuildVersion, BuildTime, CommitID)
		return
	}
	if usage {
		subcmd.PrintUsages(nil)
		return
	}
	if conf != "" {
		args = args[2:]
	}

	if err := config.InitConfig(conf); err != nil {
		panic(err)
	}
	
	if err := nameclient.InitClient(); err != nil {
		panic(err)
	}

	switch args[1] {
	case "create":
		subcmd.Create(args[2:])
	case "repo":
		subcmd.ScenRepo(args[2:])
	case "apply":
		subcmd.Apply(args[2:])
	case "get":
		subcmd.Get(args[2:])
	case "status":
		subcmd.Status(args[2:])
	case "group":
		subcmd.Group(args[2:])
	case "groups":
		subcmd.Groups(args[2:])
	case "log":
		subcmd.Log(args[2:])
	case "fix":
		subcmd.Fix(args[2:])
	case "version":
		subcmd.Version(args[2:])
	case "reset":
		subcmd.Reset(args[2:])
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
	case "xrun":
		subcmd.XRun(args[2:])
	case "pakma":
		subcmd.Pakma(args[2:])
	case "ssh":
		subcmd.SSH(args[2:])
	case "setssh":
		subcmd.SetSSH(args[2:])
	case "getssh":
		subcmd.GetSSH(args[2:])
	case "delssh":
		subcmd.DelSSH(args[2:])
	case "help":
		subcmd.PrintUsages(args[2:])
	default:
		output.PrintFatalln("sub command not support")
	}
}
