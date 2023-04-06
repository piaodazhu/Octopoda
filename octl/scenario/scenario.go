package scenario

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"octl/config"
	"os"
	"os/exec"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func ScenarioApply(file string, target string, message string) {
	buf, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	var configuration ScenarioConfigModel
	err = yaml.Unmarshal(buf, &configuration)
	if err != nil {
		panic(err)
	}

	err = checkConfig(&configuration)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch target {
		// ???
	case "prepare":
		ScenarioPrepare(&configuration, message)
		ScenarioRun(&configuration, "prepare", message)
	case "default":
		ScenarioPrepare(&configuration, message)
		ScenarioRun(&configuration, "prepare", message)
		ScenarioRun(&configuration, "start", message)
	case "purge":
		// ScenarioRun(&configuration, target, message)
		ScenarioPurge(&configuration)
	default:
		ScenarioRun(&configuration, target, message)
	}
}

func ScenarioPrepare(configuration *ScenarioConfigModel, message string) {
	var err error
	tmpName := fmt.Sprintf("%d.tar", time.Now().Nanosecond())
	
	// create this scenario
	fmt.Println(">> create scenario ", configuration.Name)
	err = ScenarioCreate(configuration.Name, configuration.Description)
	if err != nil {
		// panic(err)
	}

	// for each application
	for i := range configuration.Applications {
		app := configuration.Applications[i]
		fmt.Println(">> create app ", app.Name)

		// pack the file
		cmd := exec.Command("tar", "-cf", tmpName, "-C", app.SourcePath, ".")
		fmt.Println(cmd.String())
		err = cmd.Run()
		if err != nil {
			panic("cmd.Run")
		}
		// err = cmd.Wait()
		// if err != nil {
		// 	panic(err)
		// }
		if !cmd.ProcessState.Success() {
			panic("tar error")
		}

		// distrib the files
		f, err := os.OpenFile(tmpName, os.O_RDONLY, os.ModePerm)
		if err != nil {
			panic("err")
		}

		nodes_serialized, _ := json.Marshal(&app.Nodes)

		bodyBuffer := bytes.Buffer{}
		bodyWriter := multipart.NewWriter(&bodyBuffer)
		fileWriter, _ := bodyWriter.CreateFormFile("files", tmpName)
		io.Copy(fileWriter, f)
		bodyWriter.WriteField("appName", app.Name)
		bodyWriter.WriteField("scenario", configuration.Name)
		bodyWriter.WriteField("description", app.Description)
		bodyWriter.WriteField("message", message)

		bodyWriter.WriteField("targetNodes", string(nodes_serialized))

		contentType := bodyWriter.FormDataContentType()

		f.Close()
		bodyWriter.Close()

		url := fmt.Sprintf("http://%s:%d/%s%s",
			config.GlobalConfig.Server.Ip,
			config.GlobalConfig.Server.Port,
			config.GlobalConfig.Server.ApiPrefix,
			config.GlobalConfig.Api.ScenarioAppCreate,
		)

		res, err := http.Post(url, contentType, &bodyBuffer)
		if err != nil {
			panic("post")
		}
		defer res.Body.Close()
		msg, err := io.ReadAll(res.Body)
		if err != nil {
			panic("ReadAll")
		}

		fmt.Println(string(msg))
	}
	os.Remove(tmpName)
}

func ScenarioRun(configuration *ScenarioConfigModel, target, message string) {
	fmt.Println("- Target: ", target)
	// for each application
	for i := range configuration.Applications {
		app := configuration.Applications[i]
		fmt.Println(">> deploy app ", app.Name)
		// construct script path
		fname := ""
		for j := range app.Script {
			if app.Script[j].Target == target {
				fname = app.Script[j].File
				break
			}
		}
		if len(fname) == 0 {
			fmt.Println("invalid target")
			continue 
		}

		var sb strings.Builder
		sb.WriteString(app.ScriptPath)
		if sb.String()[sb.Len() - 1] != '/' {
			sb.WriteByte('/')
		}
		sb.WriteString(fname)
		
		// load the script
		f, err := os.OpenFile(sb.String(), os.O_RDONLY, os.ModePerm)
		if err != nil {
			panic("err")
		}

		nodes_serialized, _ := json.Marshal(&app.Nodes)

		bodyBuffer := bytes.Buffer{}
		bodyWriter := multipart.NewWriter(&bodyBuffer)
		fileWriter, _ := bodyWriter.CreateFormFile("script", fname)
		io.Copy(fileWriter, f)
		bodyWriter.WriteField("appName", app.Name)
		bodyWriter.WriteField("scenario", configuration.Name)
		bodyWriter.WriteField("description", app.Description)
		bodyWriter.WriteField("message", message)

		bodyWriter.WriteField("targetNodes", string(nodes_serialized))

		contentType := bodyWriter.FormDataContentType()

		f.Close()
		bodyWriter.Close()

		// distrib the files
		url := fmt.Sprintf("http://%s:%d/%s%s",
			config.GlobalConfig.Server.Ip,
			config.GlobalConfig.Server.Port,
			config.GlobalConfig.Server.ApiPrefix,
			config.GlobalConfig.Api.ScenarioAppDepoly,
		)

		res, err := http.Post(url, contentType, &bodyBuffer)
		if err != nil {
			panic("post")
		}
		defer res.Body.Close()
		msg, err := io.ReadAll(res.Body)
		if err != nil {
			panic("ReadAll")
		}
		fmt.Println(string(msg))
	}
}

func ScenarioPurge(configuration *ScenarioConfigModel) {
	// for each apps
	for i := range configuration.Applications {
		app := configuration.Applications[i]
		fmt.Println(">> delete app ", app.Name)

		// construct purge script path
		fname := ""
		for j := range app.Script {
			if app.Script[j].Target == "purge" {
				fname = app.Script[j].File
				break
			}
		}
		if len(fname) == 0 {
			fmt.Println("missing purge script")
			continue 
		}

		var sb strings.Builder
		sb.WriteString(app.ScriptPath)
		if sb.String()[sb.Len() - 1] != '/' {
			sb.WriteByte('/')
		}
		sb.WriteString(fname)
		
		// load the script
		f, err := os.OpenFile(sb.String(), os.O_RDONLY, os.ModePerm)
		if err != nil {
			panic("err")
		}

		nodes_serialized, _ := json.Marshal(&app.Nodes)

		bodyBuffer := bytes.Buffer{}
		bodyWriter := multipart.NewWriter(&bodyBuffer)
		fileWriter, _ := bodyWriter.CreateFormFile("script", fname)
		io.Copy(fileWriter, f)
		bodyWriter.WriteField("appName", app.Name)
		bodyWriter.WriteField("scenario", configuration.Name)

		bodyWriter.WriteField("targetNodes", string(nodes_serialized))

		contentType := bodyWriter.FormDataContentType()

		f.Close()
		bodyWriter.Close()

		// run purge script in corresponding nodes
		url := fmt.Sprintf("http://%s:%d/%s%s",
			config.GlobalConfig.Server.Ip,
			config.GlobalConfig.Server.Port,
			config.GlobalConfig.Server.ApiPrefix,
			config.GlobalConfig.Api.ScenarioAppDepoly,
		)

		res, err := http.Post(url, contentType, &bodyBuffer)
		if err != nil {
			panic("post")
		}
		defer res.Body.Close()
		msg, err := io.ReadAll(res.Body)
		if err != nil {
			panic("ReadAll")
		}
		fmt.Println(string(msg))
	}

	// purge this scenario
	err := ScenarioDelete(configuration.Name)
	if err != nil {
		fmt.Println(err)
	}
}

type ErrDupScenario struct{}
func (ErrDupScenario) Error() string { return "ErrDupScenario" }

func ScenarioCreate(name, description string) error {
	url := fmt.Sprintf("http://%s:%d/%s%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.ScenarioInfo,
	)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("name", name)
	writer.WriteField("description", description)
	writer.Close()

	res, err := http.Post(url, writer.FormDataContentType(), body)
	if err != nil {
		return err 
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return ErrDupScenario{}
	}
	return nil
}

func ScenarioDelete(name string) error {
	url := fmt.Sprintf("http://%s:%d/%s%s?name=%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.ScenarioInfo,
		name,
	)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err 
	}
	defer res.Body.Close()

	msg, _ := io.ReadAll(res.Body)
	fmt.Println(string(msg))
	return nil
}