package config

import "strings"

// special environment variables

func OctopodaEnv(scriptDir string, scriptName string, output string) []string {
	env := []string{
		"OCTOPODA_NODENAME=" + GlobalConfig.Name,
		// "OCTOPODA_NODEIP=" + GlobalConfig.Ip,
		"OCTOPODA_ROOTDIR=" + GlobalConfig.Workspace.Root + scriptDir,
		"OCTOPODA_STOREDIR=" + GlobalConfig.Workspace.Store,
		"OCTOPODA_FILENAME=" + scriptName,
		"OCTOPODA_OUTPUT=" + output,
	}
	s := strings.Split(scriptDir, "@")
	if len(s) == 2 && len(s[0]) > 0 && len(s[1]) > 0 {
		env = append(env, "OCTOPODA_APP=" + s[0])
		env = append(env, "OCTOPODA_SCENARIO=" + s[1])
	}
	return env
}