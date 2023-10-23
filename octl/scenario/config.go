package scenario

import (
	"fmt"
	"octl/node"
	"os"
	"strings"
)

type ScenarioConfigModel struct {
	Name         string                   `yaml:"name"`
	Description  string                   `yaml:"description"`
	Applications []ApplicationConfigModel `yaml:"applications"`
}

type ApplicationConfigModel struct {
	Name        string              `yaml:"name"`
	Description string              `yaml:"description"`
	ScriptPath  string              `yaml:"scriptpath"`
	SourcePath  string              `yaml:"sourcepath"`
	Nodes       []string            `yaml:"nodes"`
	Script      []ScriptConfigModel `yaml:"script"`
}

type ScriptConfigModel struct {
	Target string `yaml:"target"`
	File   string `yaml:"file"`
	Order  int    `yaml:"order"`
}

func checkConfig(config *ScenarioConfigModel) error {
	if config.Name == "" || config.Description == "" || len(config.Applications) == 0 {
		return fmt.Errorf("missing fields in deployment.yaml: scenario.name(%t) scenario.description(%t) scenario.applications(%t)",
			config.Name == "", config.Description == "", len(config.Applications) == 0)
	}

	// nodeset := map[string]struct{}{}
	for i := range config.Applications {
		app := &config.Applications[i]
		if app.Name == "" || app.Description == "" || app.ScriptPath == "" || len(app.Nodes) == 0 {
			return fmt.Errorf("missing fields in deployment.yaml: app.name(%t) app.description(%t) app.applications(%t)",
				config.Name == "", config.Description == "", len(config.Applications) == 0)
		}

		// path recorrect
		app.SourcePath = basePath + "/" + app.SourcePath
		app.ScriptPath = basePath + "/" + app.ScriptPath

		// check source path valid
		info, err := os.Stat(app.SourcePath)
		if err != nil || !info.IsDir() {
			return fmt.Errorf("invalid source filepath: %s", app.SourcePath)
		}

		app.Nodes, err = expandAlias(app.Nodes)
		if err != nil {
			return err
		}

		app.Nodes, err = node.NodesParse(app.Nodes)
		if err != nil {
			return fmt.Errorf("invalid nodes list in app %s: %v", app.Name, strings.Join(app.Nodes, ", "))
		}

		// check: must implement the 4 basic target
		err = checkTarget(app.Name, app.Script, app.ScriptPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func checkTarget(name string, script []ScriptConfigModel, path string) error {
	if len(script) < 4 {
		return fmt.Errorf("target 'prepare' 'start' 'stop' 'purge' msut be implemented in app %s", name)
	}
	mustImpl := map[string]bool{
		"prepare": false,
		"start":   false,
		"stop":    false,
		"purge":   false,
	}
	seen := map[string]struct{}{}
	for i := range script {
		target := script[i].Target
		// check duplication
		if _, found := seen[target]; found {
			return fmt.Errorf("target is already exists in app %s: %s", name, target)
		}
		seen[target] = struct{}{}

		// check mustImpl
		if _, found := mustImpl[target]; found {
			mustImpl[target] = true
		}

		// check file exists
		filePath := path + script[i].File
		info, err := os.Stat(filePath)

		if err != nil || info.IsDir() {
			return fmt.Errorf("invalid script filepath in app %s: %s", name, filePath)
		}
	}

	// check mustImpl
	for target, impl := range mustImpl {
		if !impl {
			return fmt.Errorf("target not implemented by script in app %s: %s", name, target)
		}
	}
	return nil
}

type NodeInfoText struct {
	Name         string `json:"name"`
	Health       string `json:"health"`
	MsgConnState string `json:"msg_conn"`
	OnlineTime   string `json:"online_time,omitempty"`
	OfflineTime  string `json:"offline_time,omitempty"`
	LastOnline   string `json:"last_active,omitempty"`
}

type NodesInfoText struct {
	NodeInfoList []*NodeInfoText `json:"nodes"`
	Total        int             `json:"total"`
	Active       int             `json:"active"`
	Offline      int             `json:"offline"`
}
