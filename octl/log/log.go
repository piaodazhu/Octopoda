package log

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/httpclient"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/octl/workgroup"
	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/errs"
)

func NodeLog(name string, params []string) (string, error) {
	maxlines, maxdaysbefore := 30, 0
	for i := range params {
		if len(params[i]) < 3 {
			continue
		}
		switch params[i][:2] {
		case "-l":
			x, err := strconv.Atoi(params[i][2:])
			if err != nil {
				return "invalid args", err
			}
			maxlines = x
		case "-d":
			x, err := strconv.Atoi(params[i][2:])
			if err != nil {
				return "invalid args", err
			}
			maxdaysbefore = x
		default:
		}
	}
	url := fmt.Sprintf("https://%s/%s%s?name=%s&maxlines=%d&maxdaysbefore=%d",
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodeLog,
		name,
		maxlines,
		maxdaysbefore,
	)
	req, _ := http.NewRequest("GET", url, nil)
	workgroup.SetHeader(req)
	res, err := httpclient.BrainClient.Do(req)
	if err != nil {
		emsg := "http get error."
		output.PrintFatalln(emsg, err)
		return emsg, err
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	output.PrintJSON(raw)
	return string(raw), nil
}

func PullLog(name string, maxlines, maxdaysbefore int) (*protocols.LogParams, *errs.OctlError) {
	url := fmt.Sprintf("https://%s/%s%s?name=%s&maxlines=%d&maxdaysbefore=%d",
		config.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodeLog,
		name,
		maxlines,
		maxdaysbefore,
	)
	req, _ := http.NewRequest("GET", url, nil)
	workgroup.SetHeader(req)
	res, err := httpclient.BrainClient.Do(req)
	if err != nil {
		emsg := "http get error." + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlHttpRequestError, emsg)
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		emsg := fmt.Sprintf("[%d]msg=%s\n", res.StatusCode, string(raw))
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlHttpStatusError, emsg)
	}

	logs := protocols.LogParams{}
	err = json.Unmarshal(raw, &logs)
	if err != nil {
		emsg := "res unmarshal error: " + err.Error()
		output.PrintFatalln(emsg)
		return nil, errs.New(errs.OctlMessageParseError, emsg)
	}
	output.PrintJSON(raw)
	return &logs, nil
}
