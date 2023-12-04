package scenario

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/piaodazhu/Octopoda/octl/workgroup"
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

func checkConfig(ctx context.Context, config *ScenarioConfigModel) ([]string, error) {
	logList := []string{}
	if config.Name == "" || config.Description == "" || len(config.Applications) == 0 {
		return logList, fmt.Errorf("missing fields in deployment.yaml: scenario.name(%t) scenario.description(%t) scenario.applications(%t)",
			config.Name == "", config.Description == "", len(config.Applications) == 0)
	}
	logList = append(logList, "succeed in checking basic fields of scenario")
	for i := range config.Applications {
		app := &config.Applications[i]
		if app.Name == "" || app.Description == "" || app.ScriptPath == "" || len(app.Nodes) == 0 {
			return logList, fmt.Errorf("missing fields in deployment.yaml: app.name(%t) app.description(%t) app.applications(%t)",
				config.Name == "", config.Description == "", len(config.Applications) == 0)
		}
		logList = append(logList, "succeed in checking basic fields of app "+app.Name)

		// path recorrect
		app.SourcePath = basePath + "/" + app.SourcePath
		app.ScriptPath = basePath + "/" + app.ScriptPath

		// check source path valid
		info, err := os.Stat(app.SourcePath)
		if err != nil || !info.IsDir() {
			return logList, fmt.Errorf("invalid source filepath: %s", app.SourcePath)
		}
		logList = append(logList, "succeed in checking source path of app "+app.Name)

		app.Nodes, err = expandAlias(app.Nodes)
		if err != nil {
			return logList, err
		}
		logList = append(logList, "succeed in resolving alias of app "+app.Name)

		doneChan := make(chan error, 1)
		go func() {
			app.Nodes, err = workgroup.NodesParse(app.Nodes)
			doneChan <- err
			close(doneChan)
		}()
		select {
		case err := <-doneChan:
			if err != nil {
				return logList, fmt.Errorf("invalid nodes list in app %s: %v: %s", app.Name, strings.Join(app.Nodes, ", "), err.Error())
			}
		case <-ctx.Done():
			return logList, errors.New("request canceled by context")
		}

		logList = append(logList, "succeed in parsing nodes of app "+app.Name)

		// check: must implement the 4 basic target
		err = checkTarget(app.Name, app.Script, app.ScriptPath)
		if err != nil {
			return logList, err
		}
		logList = append(logList, "succeed in checking targets of app "+app.Name)
	}
	return logList, nil
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
