package log

import (
	"fmt"
	"io"
	"net/http"
	"octl/config"
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
		switch params[i][0] {
		case 'l':
			x, err := strconv.Atoi(params[i][1:])
			if err != nil {
				return
			}
			maxlines = x
		case 'd':
			x, err := strconv.Atoi(params[i][1:])
			if err != nil {
				return
			}
			maxdaysbefore = x
		default:
			return
		}
	}
	url := fmt.Sprintf("http://%s:%d/%s%s?name=%s&maxlines=%d&maxdaysbefore=%d",
		config.GlobalConfig.Server.Ip,
		config.GlobalConfig.Server.Port,
		config.GlobalConfig.Server.ApiPrefix,
		config.GlobalConfig.Api.NodeLog,
		name,
		maxlines,
		maxdaysbefore,
	)
	res, err := http.Get(url)
	if err != nil {
		panic("Get")
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)

	// fmt.Print(string(raw))
	output.PrintJSON(raw)
}
