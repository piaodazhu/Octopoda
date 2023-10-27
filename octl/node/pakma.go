package node

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"octl/config"
	"octl/nameclient"
	"octl/output"
	"strconv"
	"time"
	"unicode"
)

const timefmt string = "2006-01-02@15:04:05"

func Pakma(firstarg string, args []string) (string, error) {
	var subcmd string = firstarg
	var version string = ""
	var names []string

	var timestr, limit string
	END := len(args)

	for i := range args {
		if len(args[i]) < 3 {
			continue
		}
		switch args[i][:2] {
		case "-t":
			timestr = args[i][2:]
			_, err := time.Parse(timefmt, timestr)
			if err != nil {
				emsg := fmt.Sprintf("pakma subcmd: invalid timestr (should be like %s).", timefmt)
				output.PrintFatalln(emsg, err)
				return emsg, err
			}
			END--
		case "-l":
			limit = args[i][2:]
			x, err := strconv.Atoi(limit)
			if err != nil || x <= 0 {
				emsg := "pakma subcmd: invalid limit (should be int >0)"
				output.PrintFatalln(emsg, err)
				return emsg, err
			}
			END--
		default:
		}
	}
	if END != len(args) && firstarg != "history" {
		emsg := "only packma history support -t and -l"
		output.PrintFatalln(emsg)
		return emsg, errors.New(emsg)
	}

	if len(args) < 1 {
		emsg := "not any node is specified"
		output.PrintFatalln(emsg)
		return emsg, errors.New(emsg)
	}

	switch firstarg {
	case "install":
		fallthrough // same process as upgrade
	case "upgrade":
		if len(args) < 2 {
			emsg := "not any node is specified"
			output.PrintFatalln(emsg)
			return emsg, errors.New(emsg)
		}
		version = args[0]
		names = args[1:]
	case "state":
		fallthrough
	case "confirm":
		fallthrough
	case "cancel":
		fallthrough
	case "clean":
		fallthrough
	case "downgrade":
		names = args
	case "history":
		names = args[:END]
	default:
		emsg := fmt.Sprintf("pakma subcommand not support: %s.", firstarg)
		output.PrintFatalln(emsg)
		return emsg, errors.New(emsg)
	}

	names_filtered := []string{}
	hasBrain := false
	for i := range names {
		if names[i] != "brain" {
			names_filtered = append(names_filtered, names[i])
		} else {
			hasBrain = true
		}
	}

	nodes, err := NodesParse(names_filtered)
	if err != nil {
		msg := "node parse."
		output.PrintFatalln(msg, err)
		return msg, err
	}

	if hasBrain {
		nodes = append(nodes, "brain")
	}

	URL := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.API_Pakma,
	)
	nodes_serialized, _ := config.Jsoner.Marshal(&nodes)

	values := url.Values{}
	values.Set("command", subcmd)
	values.Set("targetNodes", string(nodes_serialized))
	values.Set("time", timestr)
	values.Set("limit", limit)

	if version != "" {
		if checkVersion(version) {
			values.Set("version", version)
		} else {
			emsg := fmt.Sprintf("pakma invalid version number (right example: 1.4.1): %s.", version)
			output.PrintFatalln(emsg)
			return emsg, errors.New(emsg)
		}
	}

	res, err := http.PostForm(URL, values)
	if err != nil {
		emsg := "http post error."
		output.PrintFatalln(emsg, err)
		return emsg, err
	}
	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)
	output.PrintJSON(raw)
	return string(raw), err
}

func checkVersion(version string) bool {
	dotCnt := 0
	for _, c := range version {
		if c == '.' {
			dotCnt++
		} else if !unicode.IsNumber(c) {
			return false
		}
	}
	if version[0] == '.' || version[len(version)-1] == '.' {
		return false
	}
	return true
}
