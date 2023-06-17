package task

import (
	"fmt"
	"io"
	"net/http"
	"octl/config"
	"octl/httpnc"
	"time"

	"github.com/briandowns/spinner"
)

type ErrWaitTask struct {
	status  int
	message string
}

func (e ErrWaitTask) Error() string { return fmt.Sprintf("[%d] %s\n", e.status, e.message) }
func WaitTask(prefix string, taskid string) ([]byte, error) {
	url := fmt.Sprintf("http://%s/%s%s?taskid=%s",
		httpnc.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.GlobalConfig.Api.TaskState,
		taskid,
	)
	// fmt.Fprintf(os.Stdout, "PROCESSING  ")
	s := spinner.New(spinner.CharSets[7], 200*time.Millisecond)
	s.Prefix = prefix
	s.Start()
	defer s.Stop()

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
			fmt.Println("  [DONE]")
			return msg, nil
		} else if res.StatusCode == 202 {
			time.Sleep(time.Second * 1)
		} else {
			fmt.Println("  [FAILED]")
			return nil, ErrWaitTask{res.StatusCode, string(msg)}
		}
	}
}
