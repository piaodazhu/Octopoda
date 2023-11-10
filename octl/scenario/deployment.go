package scenario

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/nameclient"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/octl/shell"
	"github.com/piaodazhu/Octopoda/octl/task"
	"github.com/piaodazhu/Octopoda/protocols/errs"

	"github.com/mholt/archiver/v3"
	"gopkg.in/yaml.v3"
)

var basePath string

func ScenarioApply(ctx context.Context, scenFolder string, target string, message string) ([]string, *errs.OctlError) {
	var err error
	logList := []string{}
	basePath, err = filepath.Abs(scenFolder)
	if err != nil {
		emsg := fmt.Sprintf("scenario dir %s not exists: %s\n", scenFolder, err.Error())
		output.PrintFatalln(emsg)
		return logList, errs.New(errs.OctlFileOperationError, emsg)
	}
	logList = append(logList, "succeed in locating scenario dir")

	confFile := basePath + "/deployment.yaml"
	buf, err := os.ReadFile(confFile)
	if err != nil {
		emsg := fmt.Sprintf("os.ReadFile(%s): %s", confFile, err.Error())
		output.PrintFatalln(emsg)
		return logList, errs.New(errs.OctlReadConfigError, emsg)
	}
	logList = append(logList, "succeed in reading deployment file")

	var configuration ScenarioConfigModel
	err = yaml.Unmarshal(buf, &configuration)
	if err != nil {
		emsg := "yaml.Unmarshal(): " + err.Error()
		output.PrintFatalln(emsg)
		return logList, errs.New(errs.OctlReadConfigError, emsg)
	}
	logList = append(logList, "succeed in unmarshalling deployment file")

	aliasFile := basePath + "/alias.yaml"
	err = parseAliasFile(aliasFile)
	if err != nil {
		emsg := fmt.Sprintf("parseAliasFile(%s): %s", aliasFile, err.Error())
		output.PrintFatalln(emsg)
		return logList, errs.New(errs.OctlReadConfigError, emsg)
	}
	logList = append(logList, "succeed in parsing alias file")

	checkLogList, err := checkConfig(ctx, &configuration)
	logList = append(logList, checkLogList...)
	if err != nil {
		emsg := "checkConfig(): " + err.Error()
		output.PrintFatalln(emsg)
		return logList, errs.New(errs.OctlReadConfigError, emsg)
	}
	logList = append(logList, "succeed in checking deployment file")

	var subLogList []string
	var subError *errs.OctlError

	switch target {
	// ???
	case "prepare":
		subLogList, subError = ScenarioPrepare(ctx, &configuration, "(Prepare-Stage1) "+message)
		logList = append(logList, subLogList...)
		if err != nil {
			return logList, subError
		}

		subLogList, subError = ScenarioRun(ctx, &configuration, "prepare", "(Prepare-Stage2) "+message)
	case "default":
		subLogList, subError = ScenarioPrepare(ctx, &configuration, "(Prepare-Stage1) "+message)
		logList = append(logList, subLogList...)
		if subError != nil {
			return logList, subError
		}

		subLogList, subError = ScenarioRun(ctx, &configuration, "prepare", "(Prepare-Stage2) "+message)
		logList = append(logList, subLogList...)
		if subError != nil {
			return logList, subError
		}

		subLogList, err = ScenarioRun(ctx, &configuration, "start", "(Deploy) "+message)
	case "purge":
		subLogList, subError = ScenarioPurge(ctx, &configuration)
	default:
		subLogList, subError = ScenarioRun(ctx, &configuration, target, "(Deploy) "+message)
	}

	logList = append(logList, subLogList...)
	logList = append(logList, "done")
	return logList, subError
}

func ScenarioPrepare(ctx context.Context, configuration *ScenarioConfigModel, message string) ([]string, *errs.OctlError) {
	var logList []string
	// create this scenario
	result, err := ScenarioCreate(ctx, configuration.Name, configuration.Description)
	logList = append(logList, result...)
	if err != nil {
		return logList, err
	}

	// for each application
	for i := range configuration.Applications {
		app := configuration.Applications[i]
		info := app.Name + "@" + configuration.Name
		output.PrintInfoln(">> create", info)

		if app.SourcePath[len(app.SourcePath)-1] == '/' {
			app.SourcePath = app.SourcePath + "."
		}
		packName := fmt.Sprintf("%d.zip", time.Now().Nanosecond())
		archiver.DefaultZip.OverwriteExisting = true
		err := archiver.DefaultZip.Archive([]string{app.SourcePath}, packName)
		if err != nil {
			emsg := fmt.Sprintf("archiver.DefaultZip.Archive([]string{%s}, %s): %s", app.SourcePath, packName, err.Error())
			output.PrintFatalln(emsg)
			return logList, errs.New(errs.OctlFileOperationError, emsg)
		}
		logList = append(logList, "succeed in packing resources of app "+app.Name)

		// distrib the files
		f, err := os.OpenFile(packName, os.O_RDONLY, os.ModePerm)
		if err != nil {
			emsg := fmt.Sprintf("os.OpenFile(%s, os.O_RDONLY, os.ModePerm): %s", packName, err.Error())
			output.PrintFatalln(emsg)
			return logList, errs.New(errs.OctlFileOperationError, emsg)
		}
		logList = append(logList, "succeed in reading resources pack of app "+app.Name)

		var taskid string
		doneChan := make(chan *errs.OctlError, 1)
		go func() {
			defer close(doneChan)
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
				emsg := "http post error: " + err.Error()
				output.PrintFatalln(emsg)
				doneChan <- errs.New(errs.OctlHttpRequestError, emsg)
				return
			}

			taskid_raw, err := io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				emsg := "http read body: " + err.Error()
				output.PrintFatalln(emsg)
				doneChan <- errs.New(errs.OctlHttpRequestError, emsg)
				return
			}
			if res.StatusCode != http.StatusAccepted {
				emsg := fmt.Sprintf("http request error status=%d. ", res.StatusCode)
				output.PrintFatalln(emsg)
				doneChan <- errs.New(errs.OctlHttpStatusError, emsg)
				return
			}
			taskid = string(taskid_raw)
			doneChan <- nil
		}()
		select {
		case err := <-doneChan:
			if err != nil {
				return logList, err
			}
		case <-ctx.Done():
			return logList, errs.New(errs.OctlContextCancelError, "request canceled by context")
		}
		logList = append(logList, "succeed in distributing resources of app (get taskid) "+app.Name)

		sigChan := make(chan os.Signal, 1)
		shouldStop := false
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigChan)
		go func(tid string) {
			select {
			case _, sigCaptured := <-sigChan:
				if sigCaptured {
					shell.RunCancel(tid)
					shouldStop = true
				}
			case <-ctx.Done():
				shell.RunCancel(tid)
				shouldStop = true
			}
		}(taskid)

		results, err := task.WaitTask("PROCESSING...", string(taskid))
		if err != nil {
			emsg := "Task processing error: " + err.Error()
			output.PrintFatalln(emsg)
			return logList, errs.New(errs.OctlTaskWaitingError, emsg)
		}
		output.PrintJSON(results)
		logList = append(logList, "succeed in distributing resources of app (task finish) "+app.Name)

		if shouldStop {
			emsg := "all tasks are cancelled and exit."
			output.PrintInfoln(emsg, err)
			return logList, nil
		}
	}
	// update this scenario
	result, err = ScenarioUpdate(ctx, configuration.Name, message)
	logList = append(logList, result...)
	if err != nil {
		emsg := fmt.Sprintf("ScenarioUpdate(%s, %s).", configuration.Name, message)
		output.PrintFatalln(emsg, err)
		return logList, err
	}
	logList = append(logList, "succeed in distributing resources of scenario")
	return logList, nil
}

type orderedReq struct {
	req   *http.Request
	order int
	info  string
}

func ScenarioRun(ctx context.Context, configuration *ScenarioConfigModel, target, message string) ([]string, *errs.OctlError) {
	var logList []string
	output.PrintInfoln("\n- Target: ", target)
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
			// output.PrintInfoln("invalid target")
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
			emsg := fmt.Sprintf("os.OpenFile(%s, os.O_RDONLY, os.ModePerm): %s", sb.String(), err.Error())
			output.PrintFatalln(emsg)
			return logList, errs.New(errs.OctlFileOperationError, emsg)
		}
		logList = append(logList, "succeed in loading script for app "+app.Name)

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

		req, err := http.NewRequest("POST", url, &bodyBuffer)
		if err != nil {
			emsg := "http post error: " + err.Error()
			output.PrintFatalln(emsg)
			return logList, errs.New(errs.OctlHttpRequestError, emsg)
		}
		req.Header.Set("Content-Type", contentType)
		logList = append(logList, "succeed in making deployment request for app "+app.Name)

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
		output.PrintInfoln(">> deploy", orlist[i].info)
		var taskid string
		doneChan := make(chan *errs.OctlError, 1)
		go func() {
			defer close(doneChan)
			res, err := http.DefaultClient.Do(orlist[i].req)
			if err != nil {
				emsg := "http.DefaultClient.Do(): " + err.Error()
				output.PrintFatalln(emsg)
				doneChan <- errs.New(errs.OctlHttpRequestError, emsg)
				return
			}

			taskid_raw, err := io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				emsg := "http read body: " + err.Error()
				output.PrintFatalln(emsg)
				doneChan <- errs.New(errs.OctlHttpRequestError, emsg)
				return
			}

			if res.StatusCode != http.StatusAccepted {
				emsg := fmt.Sprintf("http request error status=%d. ", res.StatusCode)
				output.PrintFatalln(emsg)
				doneChan <- errs.New(errs.OctlHttpStatusError, emsg)
				return
			}
			taskid = string(taskid_raw)
			doneChan <- nil
		}()
		select {
		case err := <-doneChan:
			if err != nil {
				return logList, err
			}
		case <-ctx.Done():
			return logList, errs.New(errs.OctlContextCancelError, "request canceled by context")
		}
		logList = append(logList, "succeed in running"+target+"script of app (get taskid) "+orlist[i].info)

		sigChan := make(chan os.Signal, 1)
		shouldStop := false
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigChan)
		go func(tid string) {
			select {
			case _, sigCaptured := <-sigChan:
				if sigCaptured {
					shell.RunCancel(tid)
					shouldStop = true
				}
			case <-ctx.Done():
				shell.RunCancel(tid)
				shouldStop = true
			}
		}(taskid)

		results, err := task.WaitTask("PROCESSING...", string(taskid))
		if err != nil {
			emsg := "Task processing error: " + err.Error()
			output.PrintFatalln(emsg, err)
			return logList, errs.New(errs.OctlTaskWaitingError, emsg)
		}
		output.PrintJSON(results)
		logList = append(logList, "succeed in running"+target+"script of app (task finish) "+orlist[i].info)

		if shouldStop {
			emsg := "all tasks are cancelled and exit."
			output.PrintInfoln(emsg)
			return logList, nil
		}
	}

	// update this scenario
	result, err := ScenarioUpdate(ctx, configuration.Name, message)
	logList = append(logList, result...)
	if err != nil {
		emsg := fmt.Sprintf("ScenarioUpdate(%s, %s).", configuration.Name, message)
		output.PrintFatalln(emsg, err)
		return logList, err
	}
	return logList, nil
}

func ScenarioPurge(ctx context.Context, configuration *ScenarioConfigModel) ([]string, *errs.OctlError) {
	var logList []string
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
			output.PrintWarningln("missing purge script")
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
			emsg := fmt.Sprintf("os.OpenFile(%s, os.O_RDONLY, os.ModePerm): %s", sb.String(), err.Error())
			output.PrintFatalln(emsg)
			return logList, errs.New(errs.OctlTaskWaitingError, emsg)
		}
		logList = append(logList, "succeed in loading purge script for app "+app.Name)

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
			emsg := "http post error: " + err.Error()
			output.PrintFatalln(emsg)
			return logList, errs.New(errs.OctlHttpRequestError, emsg)
		}
		req.Header.Set("Content-Type", contentType)
		logList = append(logList, "succeed in making purge request for app "+app.Name)

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
		output.PrintInfoln(">> delete", orlist[i].info)
		var taskid string
		doneChan := make(chan *errs.OctlError, 1)
		go func() {
			defer close(doneChan)
			res, err := http.DefaultClient.Do(orlist[i].req)
			if err != nil {
				emsg := "http.DefaultClient.Do(): " + err.Error()
				output.PrintFatalln(emsg)
				doneChan <- errs.New(errs.OctlHttpRequestError, emsg)
				return 
			}

			msg, err := io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				emsg := "http read body: " + err.Error()
				output.PrintFatalln(emsg)
				doneChan <- errs.New(errs.OctlHttpRequestError, emsg)
				return 
			}

			if res.StatusCode != http.StatusAccepted {
				emsg := fmt.Sprintf("http request error msg=%s, status=%d. ", msg, res.StatusCode)
				output.PrintFatalln(emsg)
				doneChan <- errs.New(errs.OctlHttpStatusError, emsg)
				return 
			}
			taskid = string(msg)
			doneChan <- nil
		}()
		select {
		case err := <-doneChan:
			if err != nil {
				return logList, err
			}
		case <-ctx.Done():
			return logList, errs.New(errs.OctlContextCancelError, "request canceled by context")
		}
		logList = append(logList, "succeed in purging app (get taskid) "+orlist[i].info)

		sigChan := make(chan os.Signal, 1)
		shouldStop := false
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigChan)
		go func(tid string) {
			select {
			case _, sigCaptured := <-sigChan:
				if sigCaptured {
					shell.RunCancel(tid)
					shouldStop = true
				}
			case <-ctx.Done():
				shell.RunCancel(tid)
				shouldStop = true
			}
		}(taskid)

		results, err := task.WaitTask("PROCESSING...", taskid)
		if err != nil {
			emsg := "Task processing error: " + err.Error()
			output.PrintFatalln(emsg)
			return logList, errs.New(errs.OctlTaskWaitingError, emsg)
		}
		output.PrintJSON(results)
		logList = append(logList, "succeed in purging app (task finish) "+orlist[i].info)

		if shouldStop {
			emsg := "all tasks are cancelled and exit."
			output.PrintInfoln(emsg)
			return logList, nil
		}
	}

	// purge this scenario
	result, err := ScenarioDelete(ctx, configuration.Name)
	logList = append(logList, result...)
	if err != nil {
		emsg := fmt.Sprintf("ScenarioDelete(%s).", configuration.Name)
		output.PrintFatalln(emsg, err)
		return logList, err
	}
	return logList, nil
}

func ScenarioCreate(ctx context.Context, name, description string) ([]string, *errs.OctlError) {
	var logList []string
	output.PrintInfoln(">> create scenario", name)
	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_ScenarioInfo,
	)

	doneChan := make(chan *errs.OctlError, 1)
	go func() {
		defer close(doneChan)
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("name", name)
		writer.WriteField("description", description)
		writer.Close()

		res, err := http.Post(url, writer.FormDataContentType(), body)
		if err != nil {
			emsg := "http post error: " + err.Error()
			output.PrintFatalln(emsg)
			doneChan <- errs.New(errs.OctlHttpRequestError, emsg)
			return
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			emsg := fmt.Sprintf("scenario %s already exists, status=%d.", name, res.StatusCode)
			output.PrintFatalln(emsg)
			doneChan <- errs.New(errs.OctlHttpStatusError, emsg)
			return
		}
		doneChan <- nil
	}()
	select {
	case err := <-doneChan:
		if err != nil {
			return logList, err
		}
	case <-ctx.Done():
		return logList, errs.New(errs.OctlContextCancelError, "request canceled by context")
	}

	logList = append(logList, "succeed in creating scenario object on brain")
	return logList, nil
}

func ScenarioUpdate(ctx context.Context, name, message string) ([]string, *errs.OctlError) {
	var logList []string
	output.PrintInfoln(">> update scenario", name)
	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_ScenarioUpdate,
	)

	doneChan := make(chan *errs.OctlError, 1)
	go func() {
		defer close(doneChan)
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("name", name)
		writer.WriteField("message", message)
		writer.Close()

		res, err := http.Post(url, writer.FormDataContentType(), body)
		if err != nil {
			emsg := "http post error: " + err.Error()
			output.PrintFatalln(emsg)
			doneChan <- errs.New(errs.OctlHttpRequestError, emsg)
			return
		}
		defer res.Body.Close()

		msg, _ := io.ReadAll(res.Body)
		output.PrintJSON(msg)
		doneChan <- nil
	}()

	select {
	case err := <-doneChan:
		if err != nil {
			return logList, err
		}
	case <-ctx.Done():
		return logList, errs.New(errs.OctlContextCancelError, "request canceled by context")
	}

	logList = append(logList, "succeed in updating scenario object on brain")
	return logList, nil
}

func ScenarioDelete(ctx context.Context, name string) ([]string, *errs.OctlError) {
	var logList []string
	output.PrintInfoln(">> delete scenario", name)
	url := fmt.Sprintf("http://%s/%s%s?name=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_ScenarioInfo,
		name,
	)

	doneChan := make(chan *errs.OctlError, 1)
	go func() {
		defer close(doneChan)
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			emsg := "http new delete request error: " + err.Error()
			output.PrintFatalln(emsg)
			doneChan <- errs.New(errs.OctlHttpRequestError, emsg)
			return
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			emsg := "http.DefaultClient.Do(): " + err.Error()
			output.PrintFatalln(emsg)
			doneChan <- errs.New(errs.OctlHttpRequestError, emsg)
			return
		}
		defer res.Body.Close()

		msg, _ := io.ReadAll(res.Body)
		output.PrintJSON(msg)
		doneChan <- nil
	}()

	select {
	case err := <-doneChan:
		if err != nil {
			return logList, err
		}
	case <-ctx.Done():
		return logList, errs.New(errs.OctlContextCancelError, "request canceled by context")
	}

	logList = append(logList, "succeed in deleting scenario object on brain")

	return logList, nil
}
