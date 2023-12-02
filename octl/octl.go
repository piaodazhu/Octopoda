package main

import (
	"flag"
	"os"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/httpclient"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/octl/subcmd"
	"github.com/piaodazhu/Octopoda/octl/workgroup"
)

// Build Information
var (
	BuildVersion string = "dev"
	BuildTime    string
	BuildName    string = "brain"
	CommitID     string = "snapshot"
)

func main() {
	// use environment variable to control command line print
	if _, found := os.LookupEnv("OCTL_NOPRINT"); !found {
		output.EnablePrint()
	}
	if _, found := os.LookupEnv("OCTL_NOCOLOR"); !found {
		output.EnableColor()
	}
	if _, found := os.LookupEnv("OCTL_NOSPINNER"); !found {
		output.EnableSpinner()
	}

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
		args = args[3:]
	} else {
		args = args[1:]
	}

	if err := config.InitConfig(conf); err != nil {
		output.PrintFatalf("cannot load config file %s: %s\n", conf, err.Error())
	}

	if err := httpclient.InitClients(); err != nil {
		output.PrintFatalf("cannot init http nameclient: %s\n", err.Error())
	}

	if err := workgroup.InitWorkgroup(httpclient.BrainClient); err != nil {
		output.PrintFatalf("cannot init workgroup: %s\n", err.Error())
	}

	if len(args) == 0 {
		subcmd.PrintUsages(nil)
		return
	}
	subcmd.Execute(args[0], args[1:])
}
