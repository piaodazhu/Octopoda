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
		ScenarioPrepare(&configuration, "(Prepare-1) "+message)
		// time.Sleep(time.Minute)
		ScenarioRun(&configuration, "prepare", "(Prepare-2) "+message)
	case "default":
		ScenarioPrepare(&configuration, "(Prepare-1) "+message)
		ScenarioRun(&configuration, "prepare", "(Prepare-2) "+message)
		ScenarioRun(&configuration, "start", "(Deploy) "+message)
	case "purge":
		// ScenarioRun(&configuration, target, message)
		ScenarioPurge(&configuration)
	default:
		ScenarioRun(&configuration, target, "(Deploy) "+message)
	}
}

func ScenarioPrepare(configuration *ScenarioConfigModel, message string) {
	var err error
	tmpName := fmt.Sprintf("%d.tar", time.Now().Nanosecond())
	
	// create this scenario
	err = ScenarioCreate(configuration.Name, configuration.Description)
	if err != nil {
		panic(err)
	}

	// for each application
	for i := range configuration.Applications {
		app := configuration.Applications[i]
		fmt.Println(">> create app ", app.Name)

		// pack the file
		cmd := exec.Command("tar", "-cf", tmpName, "-C", app.SourcePath, ".")
		// fmt.Println(cmd.String())
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

		client := http.Client{Timeout: 0}
		res, err := client.Post(url, contentType, &bodyBuffer)
		if err != nil {
			panic("post")
		}
		// defer res.Body.Close()
		msg, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			panic("ReadAll")
		}

		if res.StatusCode != 202 {
			fmt.Println("Request submit error: ", string(msg))
			return
		}
		results, err := waitTask(string(msg))
		if err != nil {
			fmt.Println("Task processing error: ", err.Error())
			return
		}
		fmt.Println(results)
	}
	os.Remove(tmpName)
	// update this scenario
	err = ScenarioUpdate(configuration.Name, message)
	if err != nil {
		fmt.Println(err)
	}
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
			// fmt.Println("invalid target")
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
		msg, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			panic("ReadAll")
		}

		if res.StatusCode != 202 {
			fmt.Println("Request submit error: ", string(msg))
			return
		}
		results, err := waitTask(string(msg))
		if err != nil {
			fmt.Println("Task processing error: ", err.Error())
			return
		}
		fmt.Println(results)
	}
	// update this scenario
	err := ScenarioUpdate(configuration.Name, message)
	if err != nil {
		fmt.Println(err)
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

		// pesudo field
		bodyWriter.WriteField("description", "PURGED")
		bodyWriter.WriteField("message", "PURGED")

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
		msg, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			panic("ReadAll")
		}

		if res.StatusCode != 202 {
			fmt.Println("Request submit error: ", string(msg))
			return
		}
		results, err := waitTask(string(msg))
		if err != nil {
			fmt.Println("Task processing error: ", err.Error())
			return
		}
		fmt.Println(results)
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
	fmt.Println(">> create scenario", name)
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

func ScenarioUpdate(name, message string) error {
	fmt.Println(">> update scenario", name)
	url := fmt.Sprintf("http://%s:%d/%s%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.ScenarioUpdate,
	)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("name", name)
	writer.WriteField("message", message)
	writer.Close()

	res, err := http.Post(url, writer.FormDataContentType(), body)
	if err != nil {
		return err 
	}
	defer res.Body.Close()

	msg, _ := io.ReadAll(res.Body)
	fmt.Println(string(msg))
	return nil
}

func ScenarioDelete(name string) error {
	fmt.Println(">> delete scenario", name)
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

type ErrWaitTask struct {
	status int
	message string
}
func (e ErrWaitTask) Error() string { return fmt.Sprintf("[%d] %s\n", e.status, e.message) }
func waitTask(taskid string) (string, error) {
	url := fmt.Sprintf("http://%s:%d/%s%s?taskid=%s",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.TaskState,
		taskid,
	)
	fmt.Print("PROCESSING...")
	time.Sleep(time.Millisecond * 100)
	for {
		res, err := http.Get(url)
		if err != nil {
			return "", err
		}
		
		msg, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return "", err
		}
		
		if res.StatusCode == 200 {
			fmt.Println("DONE")
			return string(msg), nil
		} else if res.StatusCode == 202 {
			fmt.Print(".")
			time.Sleep(time.Second * 1)
		} else {
			fmt.Println("FAILED")
			return "", ErrWaitTask{res.StatusCode, string(msg)}
		}
	}
}