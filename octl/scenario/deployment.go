package scenario

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"octl/config"
	"octl/output"
	"octl/task"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/mholt/archiver/v3"
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
		ScenarioPrepare(&configuration, "(Prepare-Stage1) "+message)
		// time.Sleep(time.Minute)
		ScenarioRun(&configuration, "prepare", "(Prepare-Stage2) "+message)
	case "default":
		ScenarioPrepare(&configuration, "(Prepare-Stage1) "+message)
		ScenarioRun(&configuration, "prepare", "(Prepare-Stage2) "+message)
		ScenarioRun(&configuration, "start", "(Deploy) "+message)
	case "purge":
		ScenarioPurge(&configuration)
	default:
		ScenarioRun(&configuration, target, "(Deploy) "+message)
	}
}

func ScenarioPrepare(configuration *ScenarioConfigModel, message string) {
	var err error
	packName := fmt.Sprintf("%d.zip", time.Now().Nanosecond())

	// create this scenario
	err = ScenarioCreate(configuration.Name, configuration.Description)
	if err != nil {
		panic(err)
	}

	// for each application
	for i := range configuration.Applications {
		app := configuration.Applications[i]
		info := app.Name + "@" + configuration.Name
		fmt.Println(">> create", info)

		// pack the file
		// cmd := exec.Command("tar", "-cf", tmpName, "-C", app.SourcePath, ".")
		// // fmt.Println(cmd.String())
		// err = cmd.Run()
		// if err != nil {
		// 	panic("cmd.Run")
		// }

		err = archiver.Archive([]string{app.SourcePath}, packName)
		if err != nil {
			panic("Archive")
		}
		// err = cmd.Wait()
		// if err != nil {
		// 	panic(err)
		// }
		// if !cmd.ProcessState.Success() {
		// 	panic("tar error")
		// }

		// distrib the files
		f, err := os.OpenFile(packName, os.O_RDONLY, os.ModePerm)
		if err != nil {
			panic("err")
		}

		nodes_serialized, _ := json.Marshal(&app.Nodes)

		bodyBuffer := bytes.Buffer{}
		bodyWriter := multipart.NewWriter(&bodyBuffer)
		fileWriter, _ := bodyWriter.CreateFormFile("files", packName)
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
		results, err := task.WaitTask("PROCESSING...", string(msg))
		if err != nil {
			fmt.Println("Task processing error: ", err.Error())
			return
		}
		output.PrintJSON(results)
	}
	os.Remove(packName)
	// update this scenario
	err = ScenarioUpdate(configuration.Name, message)
	if err != nil {
		fmt.Println(err)
	}
}

type orderedReq struct {
	req   *http.Request
	order int
	info  string
}

func ScenarioRun(configuration *ScenarioConfigModel, target, message string) {
	fmt.Println("\n- Target: ", target)
	orlist := []orderedReq{}
	// for each application
	for i := range configuration.Applications {
		app := configuration.Applications[i]
		// construct script path
		fname := ""
		order := 0
		for j := range app.Script {
			if app.Script[j].Target == target {
				fname = app.Script[j].File
				order = app.Script[j].Order
				break
			}
		}
		if len(fname) == 0 {
			// fmt.Println("invalid target")
			continue
		}

		var sb strings.Builder
		sb.WriteString(app.ScriptPath)
		if sb.String()[sb.Len()-1] != '/' {
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

		// run the scripts
		url := fmt.Sprintf("http://%s:%d/%s%s",
			config.GlobalConfig.Server.Ip,
			config.GlobalConfig.Server.Port,
			config.GlobalConfig.Server.ApiPrefix,
			config.GlobalConfig.Api.ScenarioAppDepoly,
		)

		// res, err := http.Post(url, contentType, &bodyBuffer)

		req, err := http.NewRequest("POST", url, &bodyBuffer)
		if err != nil {
			panic("NewRequest")
		}
		req.Header.Set("Content-Type", contentType)

		// add to ordered request list
		orlist = append(orlist, orderedReq{
			req:   req,
			order: order,
			info:  app.Name + "@" + configuration.Name + ": <" + target + ">",
		})
	}

	// sort the request by given order
	sort.Slice(orlist, func(i, j int) bool {
		if orlist[i].order == orlist[j].order {
			return time.Now().Minute()&1 == 1
		}
		return orlist[i].order < orlist[j].order
	})

	// perform the request one by one
	for i := range orlist {
		fmt.Println(">> deploy", orlist[i].info)
		res, err := http.DefaultClient.Do(orlist[i].req)
		if err != nil {
			panic("DoRequest")
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
		results, err := task.WaitTask("PROCESSING...", string(msg))
		if err != nil {
			fmt.Println("Task processing error: ", err.Error())
			return
		}
		output.PrintJSON(results)
	}

	// update this scenario
	err := ScenarioUpdate(configuration.Name, message)
	if err != nil {
		fmt.Println(err)
	}
}

func ScenarioPurge(configuration *ScenarioConfigModel) {
	orlist := []orderedReq{}
	// for each apps
	for i := range configuration.Applications {
		app := configuration.Applications[i]

		// construct purge script path
		fname := ""
		order := 0
		for j := range app.Script {
			if app.Script[j].Target == "purge" {
				fname = app.Script[j].File
				order = app.Script[j].Order
				break
			}
		}
		if len(fname) == 0 {
			fmt.Println("missing purge script")
			continue
		}

		var sb strings.Builder
		sb.WriteString(app.ScriptPath)
		if sb.String()[sb.Len()-1] != '/' {
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

		req, err := http.NewRequest("POST", url, &bodyBuffer)
		if err != nil {
			panic("NewRequest")
		}
		req.Header.Set("Content-Type", contentType)

		// add to ordered request list
		orlist = append(orlist, orderedReq{
			req:   req,
			order: order,
			info:  app.Name + "@" + configuration.Name + ": <purge>",
		})
	}

	// sort the request by given order
	sort.Slice(orlist, func(i, j int) bool {
		if orlist[i].order == orlist[j].order {
			return time.Now().Minute()&1 == 1
		}
		return orlist[i].order < orlist[j].order
	})

	// perform the request one by one
	for i := range orlist {
		fmt.Println(">> delete", orlist[i].info)
		res, err := http.DefaultClient.Do(orlist[i].req)
		if err != nil {
			panic("DoRequest")
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
		results, err := task.WaitTask("PROCESSING...", string(msg))
		if err != nil {
			fmt.Println("Task processing error: ", err.Error())
			return
		}
		output.PrintJSON(results)
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
	output.PrintJSON(msg)

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
	output.PrintJSON(msg)

	return nil
}
