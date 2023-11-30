package config

import "strings"

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
