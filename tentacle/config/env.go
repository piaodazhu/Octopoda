package config

import (
	"fmt"
	"strings"
)

// special environment variables

func OctopodaEnv(scriptDir string, scriptName string, output string) []string {
	env := []string{
		"OCTOPODA_NODENAME=" + GlobalConfig.Name,
		"OCTOPODA_CURRENTDIR=" + scriptDir,
		"OCTOPODA_STOREDIR=" + GlobalConfig.Workspace.Store,
		"OCTOPODA_FILENAME=" + scriptName,
		"OCTOPODA_OUTPUT=" + output,
	}
	s := strings.Split(scriptDir, "@")
	if len(s) == 2 && len(s[0]) > 0 && len(s[1]) > 0 {
		env = append(env, "OCTOPODA_APP="+s[0])
		env = append(env, "OCTOPODA_SCENARIO="+s[1])
	}
	for _, customized := range GlobalConfig.CustomEnv {
		env = append(env, fmt.Sprintf("%s=%s", customized.Key, customized.Value))
	}
	return env
}

var pathEnv map[string]string
func initPathEnv() {
	pathEnv = make(map[string]string)
	pathEnv["@root"] = "/root/"
	pathEnv["@workspace"] = GlobalConfig.Workspace.Root
	pathEnv["@fstore"] = GlobalConfig.Workspace.Store
	pathEnv["@log"] = GlobalConfig.Logger.Path
	pathEnv["@pakma"] = GlobalConfig.PakmaServer.Root
}

func ParsePathWithEnv(path string) string {
	if !strings.HasPrefix(path, "@") {
		return path 
	}
	substr := strings.SplitN(path, "/", 2)
	path = pathEnv[substr[0]]
	if len(substr) == 2 {
		path += substr[1]
	}
	return path
}
