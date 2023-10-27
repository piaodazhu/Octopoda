package log

import (
	"fmt"
	"io"
	"net/http"
	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/nameclient"
	"github.com/piaodazhu/Octopoda/octl/output"
	"strconv"
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
	url := fmt.Sprintf("http://%s/%s%s?name=%s&maxlines=%d&maxdaysbefore=%d",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_NodeLog,
		name,
		maxlines,
		maxdaysbefore,
	)
	res, err := http.Get(url)
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
