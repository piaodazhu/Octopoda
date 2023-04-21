package scenario

import (
	"fmt"
	"io"
	"net/http"
	"octl/config"
	"os"
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

type ErrInvalidNode struct{}

func (ErrInvalidNode) Error() string { return "ErrInvalidNode" }

type ErrInvalidScript struct{}

func (ErrInvalidScript) Error() string { return "ErrInvalidScript" }

type ErrMissingScript struct{}

func (ErrMissingScript) Error() string { return "ErrMissingScript" }

type ErrInvalidSource struct{}

func (ErrInvalidSource) Error() string { return "ErrInvalidSource" }

type ErrMissingFields struct{}

func (ErrMissingFields) Error() string { return "ErrMissingFields" }

type ErrDupTargets struct{}

func (ErrDupTargets) Error() string { return "ErrDupTargets" }

func checkConfig(config *ScenarioConfigModel) error {
	if config.Name == "" || config.Description == "" || len(config.Applications) == 0 {
		return ErrMissingFields{}
	}
	// fmt.Printf("%+v\n", *config)
	// os.Exit(0)

	nodeset := map[string]struct{}{}
	var err error
	for _, app := range config.Applications {
		if app.Name == "" || app.Description == "" || app.ScriptPath == "" || len(app.Nodes) == 0 {
			return ErrMissingFields{}
		}

		// check source path valid
		if len(app.SourcePath) != 0 {
			info, err := os.Stat(app.SourcePath)
			if err != nil || !info.IsDir() {
				return ErrInvalidSource{}
			}
		}

		// collect all nodename then check once
		for _, node := range app.Nodes {
			nodeset[node] = struct{}{}
		}

		// check: must implement the 4 basic target
		err = checkTarget(app.Script, app.ScriptPath)
		if err != nil {
			return err
		}
	}

	// check node validity
	if !checkNodes(nodeset) {
		return ErrInvalidNode{}
	}

	return nil
}

func checkTarget(script []ScriptConfigModel, path string) error {
	if len(script) < 4 {
		return ErrMissingScript{}
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
			return ErrDupTargets{}
		}
		seen[target] = struct{}{}

		// check mustImpl
		if _, found := mustImpl[target]; found {
			mustImpl[target] = true
		}

		// check file exists
		info, err := os.Stat(path + script[i].File)
		// fmt.Println(path + script[i].File, path)
		if err != nil || info.IsDir() {
			return ErrInvalidScript{}
		}
	}

	// check mustImpl
	for script, impl := range mustImpl {
		if !impl {
			fmt.Println("script", script, "not found")
			return ErrMissingScript{}
		}
	}
	return nil
}

type NodeInfo struct {
	Id        uint32
	Name      string
	Addr      string
	State     int32
	OnlineTs  int64
	OfflineTs int64
}

func checkNodes(nodeset map[string]struct{}) bool {
	// get all nodes in the cluster
	url := fmt.Sprintf("http://%s:%d/%s%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.NodesInfo,
	)
	res, err := http.Get(url)
	if err != nil {
		panic("Get")
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	nodes := []NodeInfo{}
	err = config.Jsoner.Unmarshal(raw, &nodes)
	if err != nil {
		panic(err)
	}
	// fmt.Println(nodes)
	// fmt.Println(nodes)
	// put them into a set
	nodemap := map[string]struct{}{}
	for i := range nodes {
		if nodes[i].State == 0 {
			nodemap[nodes[i].Name] = struct{}{}
		}
	}

	// check nodeset: is all nodename in that set?
	for node := range nodeset {
		if _, ok := nodemap[node]; !ok {
			return false
		}
	}
	return true
}
