package service

import (
	"nworkerd/logger"
	"os/exec"
)

func Reboot() {
	logger.Client.Println("Reboot.")
	_, err := exec.Command("reboot").CombinedOutput()
	if err != nil {
		logger.Client.Fatal(err)
	}
}
