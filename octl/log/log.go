package log

import (
	"fmt"
	"io"
	"net/http"
	"octl/config"
	"octl/nameclient"
	"octl/output"
	"strconv"
)

// type LogParams struct {
// 	MaxLines      int
// 	MaxDaysBefore int
// 	Logs          []string
// }

func NodeLog(name string, params []string) {
	maxlines, maxdaysbefore := 30, 0
	for i := range params {
		if len(params[i]) < 3 {
			continue
		}
		switch params[i][:2] {
		case "-l":
			x, err := strconv.Atoi(params[i][2:])
			if err != nil {
				return
			}
			maxlines = x
		case "-d":
			x, err := strconv.Atoi(params[i][2:])
			if err != nil {
				return
			}
			maxdaysbefore = x
		default:
		}
	}
	url := fmt.Sprintf("http://%s/%s%s?name=%s&maxlines=%d&maxdaysbefore=%d",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.GlobalConfig.Api.NodeLog,
		name,
		maxlines,
		maxdaysbefore,
	)
	res, err := http.Get(url)
	if err != nil {
		output.PrintFatalln("Get")
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	// fmt.Print(string(raw))
	output.PrintJSON(raw)
}
