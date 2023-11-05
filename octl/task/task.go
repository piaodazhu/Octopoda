package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/nameclient"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/protocols"

	"github.com/briandowns/spinner"
)

func WaitTask(prefix string, taskid string) ([]protocols.ExecutionResults, error) {
	url := fmt.Sprintf("http://%s/%s%s?taskid=%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_TaskState,
		taskid,
	)
	if output.IsSpinnerEnabled() {
		s := spinner.New(spinner.CharSets[7], 200*time.Millisecond)
		s.Prefix = prefix
		s.Start()
		defer s.Stop()
	}

	pollCnt := 0

	for {
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}

		msg, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return nil, err
		}

		if res.StatusCode == http.StatusOK {
			if prefix != "" {
				output.PrintInfoln("  [DONE]")
			}

			results := []protocols.ExecutionResults{}
			err = json.Unmarshal(msg, &results)
			if err != nil {
				emsg := "res unmarshal error: " + err.Error()
				output.PrintFatalln(emsg)
				return nil, errors.New(emsg)
			}
			makeCompatible(results)
			return results, nil
		} else if res.StatusCode == http.StatusAccepted {
			mul := 32
			if pollCnt < 5 {
				mul = 1 << pollCnt
			}
			time.Sleep(time.Millisecond * time.Duration(mul) * 30)
			pollCnt++
		} else {
			if prefix != "" {
				output.PrintInfoln("  [FAILED]")
			}
			return nil, fmt.Errorf("wait task error. http statuscode=%d, response=%s", res.StatusCode, string(msg))
		}
	}
}

func makeCompatible(results []protocols.ExecutionResults) {
	for i := range results {
		if results[i].Code == 0 {
			results[i].ResultCompatible = "[OK]"
		} else {
			results[i].ResultCompatible = "[ERR]"
		}
	}
}