package task

import (
	"fmt"
	"io"
	"net/http"
	"octl/config"
	"octl/nameclient"
	"octl/output"
	"time"

	"github.com/briandowns/spinner"
)

func WaitTask(prefix string, taskid string) ([]byte, error) {
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

	time.Sleep(time.Millisecond * 100)
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

		if res.StatusCode == 200 {
			if prefix != "" {
				output.PrintInfoln("  [DONE]")
			}
			return msg, nil
		} else if res.StatusCode == 202 {
			time.Sleep(time.Second * 1)
		} else {
			if prefix != "" {
				output.PrintInfoln("  [FAILED]")
			}
			return nil, fmt.Errorf("wait task error. http statuscode=%d, response=%s", res.StatusCode, string(msg))
		}
	}
}
