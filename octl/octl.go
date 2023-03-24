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
		fmt.Println("hello world")
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
	case "purge":
		subcmd.Purge(args[2:])
	default:
		fmt.Println("not support")
	}
}
