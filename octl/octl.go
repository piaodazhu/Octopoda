package main

import (
	"fmt"
	"octl/config"
	"octl/subcmd"
	"os"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Octopoda Controlling Tool. ©2023-2023 Z. Luo. All Rights Reserved.")
		return
	}
	config.InitConfig()

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
	case "help":
		subcmd.PrintUsages()
	default:
		fmt.Println("sub command not support")
	}
}
