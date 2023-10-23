package scenario

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"octl/config"
	"octl/nameclient"
	"octl/output"
	"octl/shell"
	"octl/task"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/mholt/archiver/v3"
	"gopkg.in/yaml.v3"
)

var basePath string

func ScenarioApply(scenFolder string, target string, message string) {
	var err error
	basePath, err = filepath.Abs(scenFolder)
	if err != nil {
		output.PrintFatalf("scenario dir %s not exists: %s\n", scenFolder, err.Error())
	}

	confFile := basePath + "/deployment.yaml"
	buf, err := os.ReadFile(confFile)
	if err != nil {
		output.PrintFatalln(err.Error())
	}
	var configuration ScenarioConfigModel
	err = yaml.Unmarshal(buf, &configuration)
	if err != nil {
		output.PrintFatalln(err.Error())
	}

	aliasFile := basePath + "/alias.yaml"
	err = parseAliasFile(aliasFile)
	if err != nil {
		output.PrintFatalln(err.Error())
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
	// create this scenario
	err = ScenarioCreate(configuration.Name, configuration.Description)
	if err != nil {
		output.PrintFatalln(err.Error())
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
		// 	output.PrintFatalln("cmd.Run")
		// }
		if app.SourcePath[len(app.SourcePath)-1] == '/' {
			app.SourcePath = app.SourcePath + "."
		}
		packName := fmt.Sprintf("%d.zip", time.Now().Nanosecond())
		archiver.DefaultZip.OverwriteExisting = true
		err = archiver.DefaultZip.Archive([]string{app.SourcePath}, packName)
		if err != nil {
			output.PrintFatalln("Archive")
		}
		// err = cmd.Wait()
		// if err != nil {
		// 	output.PrintFatalln(err)
		// }
		// if !cmd.ProcessState.Success() {
		// 	output.PrintFatalln("tar error")
		// }

		// distrib the files
		f, err := os.OpenFile(packName, os.O_RDONLY, os.ModePerm)
		if err != nil {
			output.PrintFatalln("err")
		}

		nodes_serialized, _ := config.Jsoner.Marshal(&app.Nodes)

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
		os.Remove(packName)
		bodyWriter.Close()

		url := fmt.Sprintf("http://%s/%s%s",
			nameclient.BrainAddr,
			config.GlobalConfig.Brain.ApiPrefix,
			config.API_ScenarioAppCreate,
		)

		client := http.Client{Timeout: 0}
		res, err := client.Post(url, contentType, &bodyBuffer)
		if err != nil {
			output.PrintFatalln("post")
		}
		// defer res.Body.Close()
		taskid, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			output.PrintFatalln("ReadAll")
		}

		sigChan := make(chan os.Signal, 1)
		shouldStop := false
		go func(tid string) {
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			_, sigCaptured := <-sigChan
			if sigCaptured {
				shell.RunCancel(tid)
				shouldStop = true
			}
		}(string(taskid))

		if res.StatusCode != 202 {
			output.PrintFatalln("Request submit error: " + string(taskid))
			return
		}
		results, err := task.WaitTask("PROCESSING...", string(taskid))
		if err != nil {
			output.PrintFatalln("Task processing error: " + err.Error())
			return
		}
		output.PrintJSON(results)

		signal.Stop(sigChan)
		if shouldStop {
			output.PrintFatalln("cancel and exit")
		}
	}
	// update this scenario
	err = ScenarioUpdate(configuration.Name, message)
	if err != nil {
		output.PrintFatalln(err.Error())
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
			output.PrintFatalln("err")
		}

		nodes_serialized, _ := config.Jsoner.Marshal(&app.Nodes)

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
		url := fmt.Sprintf("http://%s/%s%s",
			nameclient.BrainAddr,
			config.GlobalConfig.Brain.ApiPrefix,
			config.API_ScenarioAppDeploy,
		)

		// res, err := http.Post(url, contentType, &bodyBuffer)

		req, err := http.NewRequest("POST", url, &bodyBuffer)
		if err != nil {
			output.PrintFatalln("NewRequest")
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
			output.PrintFatalln("DoRequest")
		}

		taskid, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			output.PrintFatalln("ReadAll")
		}

		sigChan := make(chan os.Signal, 1)
		shouldStop := false
		go func(tid string) {
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			_, sigCaptured := <-sigChan
			if sigCaptured {
				shell.RunCancel(tid)
				shouldStop = true
			}
		}(string(taskid))

		if res.StatusCode != 202 {
			fmt.Println("Request submit error: ", string(taskid))
			return
		}
		results, err := task.WaitTask("PROCESSING...", string(taskid))
		if err != nil {
			fmt.Println("Task processing error: ", err.Error())
			return
		}
		output.PrintJSON(results)

		signal.Stop(sigChan)
		if shouldStop {
			output.PrintFatalln("cancel and exit")
		}
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
			output.PrintFatalln(err)
		}

		nodes_serialized, _ := config.Jsoner.Marshal(&app.Nodes)

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
		url := fmt.Sprintf("http://%s/%s%s",
			nameclient.BrainAddr,
			config.GlobalConfig.Brain.ApiPrefix,
			config.API_ScenarioAppDeploy,
		)

		req, err := http.NewRequest("POST", url, &bodyBuffer)
		if err != nil {
			output.PrintFatalln("NewRequest")
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
			output.PrintFatalln("DoRequest")
		}

		msg, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			output.PrintFatalln("ReadAll")
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

func ScenarioCreate(name, description string) error {
	fmt.Println(">> create scenario", name)
	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_ScenarioInfo,
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
		return fmt.Errorf("scenario %s already exists", name)
	}
	return nil
}

func ScenarioUpdate(name, message string) error {
	fmt.Println(">> update scenario", name)
	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_ScenarioUpdate,
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
	url := fmt.Sprintf("http://%s/%s%s?name=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_ScenarioInfo,
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
