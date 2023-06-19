package main

import (
	"fmt"
	"os"
	"os/exec"
	"pakma/config"
	"pakma/httpsclient"
	"sync"
	"time"
)

var PakmaError error
var Busy sync.Mutex
var State PakmaState
var PreviewTimer *time.Timer

func InitFiniteStateMachine() {
	StateMsg = map[int]string{}
	StateMsg[EMPTY] = "EMPTY"
	StateMsg[STABLE] = "STABLE"
	StateMsg[PREVIEW] = "PREVIEW"

	PakmaError = nil

	if loadState() != nil {
		initState()
	}

	if !dirExists(config.GlobalConfig.Packma.Root) {
		os.MkdirAll(config.GlobalConfig.Packma.Root, os.ModePerm)
	}
}

func doInstall(version string) {
	Busy.Lock()
	defer Busy.Unlock()

	PakmaError = installVersion(version)
	if PakmaError != nil {
		WriteHistory("Fail to install version %s as stable version", version)
		return
	}
	WriteHistory("Install version %s as stable version", version)
	State = PakmaState{
		StateType: STABLE,
		Version1:  "",
		Version2:  version,
		Version3:  "",
	}
	dumpState()
}

func doUpgrade(version string) {
	Busy.Lock()
	defer Busy.Unlock()

	PakmaError = installVersion(version)
	if PakmaError != nil {
		WriteHistory("Fail to upgrade stable version %s to version %s as preview version", State.Version2, version)
		return
	}
	// 6 change state
	WriteHistory("Upgrade stable version %s to version %s as preview version", State.Version2, version)
	State = PakmaState{
		StateType: PREVIEW,
		Version1:  State.Version1,
		Version2:  State.Version2,
		Version3:  version, // preview
	}
	dumpState()

	// auto reset
	PreviewTimer = time.AfterFunc(time.Second*time.Duration(config.GlobalConfig.Packma.PreviewDuration), doCancel)
}

func doDownGrade() {
	Busy.Lock()
	defer Busy.Unlock()

	if State.StateType != STABLE {
		return
	}

	lastVersion := State.Version1
	if lastVersion == "" {
		PakmaError = fmt.Errorf("cannot down grade because I didn't remember last version")
		return
	}

	PakmaError = installVersion(lastVersion)
	if PakmaError != nil {
		WriteHistory("Fail to downgrade version %s to version %s as stable version", State.Version2, State.Version1)
		return
	}

	WriteHistory("version %s downgrade to version %s as stable version", State.Version2, State.Version1)
	State = PakmaState{
		StateType: STABLE,
		Version1:  "",
		Version2:  lastVersion,
		Version3:  "", // preview
	}
	dumpState()
}

func doCancel() {
	Busy.Lock()
	defer Busy.Unlock()
	if PreviewTimer != nil && !PreviewTimer.Stop() {
		select {
		case <-PreviewTimer.C:
		default:
		}
	}

	PakmaError = installVersion(State.Version2)
	if PakmaError != nil {
		WriteHistory("Fail to cancel preview version %s and back to stable version %s", State.Version3, State.Version2)
		return
	}

	WriteHistory("Cancel preview version %s and back to stable version %s", State.Version3, State.Version2)
	State = PakmaState{
		StateType: STABLE,
		Version1:  State.Version1,
		Version2:  State.Version2,
		Version3:  "", // preview
	}

	dumpState()
}

func doConfirm() {
	Busy.Lock()
	defer Busy.Unlock()
	if PreviewTimer != nil && !PreviewTimer.Stop() {
		select {
		case <-PreviewTimer.C:
		default:
		}
	}

	WriteHistory("Confirm preview version %s as stable version", State.Version3)
	State = PakmaState{
		StateType: STABLE,
		Version1:  State.Version2,
		Version2:  State.Version3,
		Version3:  "", // preview
	}
	dumpState()
}

// func installVersion(version string) error {
// 	fmt.Printf("intall %s_v%s_%s_%s\n", config.GlobalConfig.AppName, version, config.GlobalConfig.AppOS, config.GlobalConfig.AppArch)

// 	return nil
// }

func installVersion(version string) error {
	// 1 check config file: config file must be installed
	if !fileExists(fmt.Sprintf("/etc/octopoda/%s/%s.yaml", config.GlobalConfig.AppName, config.GlobalConfig.AppName)) {
		return fmt.Errorf("%s not found", fmt.Sprintf("/etc/octopoda/%s/%s.yaml", config.GlobalConfig.AppName, config.GlobalConfig.AppName))
	}

	// 2 check service unit file: service unit file must be installed
	if !fileExists(fmt.Sprintf("/etc/systemd/system/%s.service", config.GlobalConfig.AppName)) {
		return fmt.Errorf("%s not found", fmt.Sprintf("/etc/systemd/system/%s.service", config.GlobalConfig.AppName))
	}

	// 3 search release package in local path
	path := getPathFromVersion(version)
	var bin string = config.GlobalConfig.Packma.Root + path + "/" + config.GlobalConfig.AppName
	var cmd *exec.Cmd
	var err error

	if fileExists(bin) {
		goto begininstall
	}

	// 4 fetch release package from httpsNameServer, and unpack
	if !fileExists(config.GlobalConfig.Packma.Root + path + ".tar.xz") {
		err := httpsclient.FetchReleasePackage(path, config.GlobalConfig.Packma.Root)
		if err != nil {
			return fmt.Errorf("cannot fetch %s.tar.xz from httpsNameServer", path)
		}
	}
	cmd = exec.Command("tar", "-Jxvf", config.GlobalConfig.Packma.Root+path+".tar.xz", "-C", config.GlobalConfig.Packma.Root)
	err = cmd.Run()
	// fmt.Println(cmd.String())
	if err != nil {
		return fmt.Errorf("cannot unpack %s.tar.xz", config.GlobalConfig.Packma.Root+path+".tar.xz")
	}

	
begininstall:
	// 5 before stopping the old service, sleep for a while to make sure tentacle/brain report its result to octl
	time.Sleep(time.Second)

	// 6 stop service && replace binary && start service
	stopCmd := exec.Command("systemctl", "stop", config.GlobalConfig.AppName)
	replaceCmd := exec.Command("cp", bin, "/usr/local/bin/octopoda/")
	startCmd := exec.Command("systemctl", "start", config.GlobalConfig.AppName)
	err = stopCmd.Run()
	if err != nil {
		return fmt.Errorf("cannot stop the old running binary")
	}
	err = replaceCmd.Run()
	if err != nil {
		return fmt.Errorf("cannot replace the old running binary to new one")
	}
	err = startCmd.Run()
	if err != nil {
		return fmt.Errorf("cannot start the new running binary")
	}
	// fmt.Println(stopCmd.String())
	// fmt.Println(replaceCmd.String())
	// fmt.Println(startCmd.String())
	return nil
}
