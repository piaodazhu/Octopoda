package node

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"unicode"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/nameclient"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/protocols/errs"
)

const timefmt string = "2006-01-02@15:04:05"

func Pakma(firstarg string, args []string) (string, *errs.OctlError) {
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
				emsg := fmt.Sprintf("pakma subcmd: invalid timestr (should be like %s): %s", timefmt, err.Error())
				output.PrintFatalln(emsg)
				return emsg, errs.New(errs.OctlArgumentError, emsg)
			}
			END--
		case "-l":
			limit = args[i][2:]
			x, err := strconv.Atoi(limit)
			if err != nil || x <= 0 {
				emsg := "pakma subcmd: invalid limit (should be int >0): " + err.Error()
				output.PrintFatalln(emsg)
				return emsg, errs.New(errs.OctlArgumentError, emsg)
			}
			END--
		default:
		}
	}
	if END != len(args) && firstarg != "history" {
		emsg := "only packma history support -t and -l"
		output.PrintFatalln(emsg)
		return emsg, errs.New(errs.OctlArgumentError, emsg)
	}

	if len(args) < 1 {
		emsg := "not any node is specified"
		output.PrintFatalln(emsg)
		return emsg, errs.New(errs.OctlArgumentError, emsg)
	}

	switch firstarg {
	case "install":
		fallthrough // same process as upgrade
	case "upgrade":
		if len(args) < 2 {
			emsg := "not any node is specified"
			output.PrintFatalln(emsg)
			return emsg, errs.New(errs.OctlArgumentError, emsg)
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
		return emsg, errs.New(errs.OctlArgumentError, emsg)
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
		emsg := "node parse error: " + err.Error()
		output.PrintFatalln(emsg)
		return "", errs.New(errs.OctlNodeParseError, emsg)
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
			return emsg, errs.New(errs.OctlArgumentError, emsg)
		}
	}

	res, err := http.PostForm(URL, values)
	if err != nil {
		emsg := "http post error: " + err.Error()
		output.PrintFatalln(emsg)
		return emsg, errs.New(errs.OctlHttpRequestError, emsg)
	}
	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)
	output.PrintJSON(raw)
	return string(raw), nil
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
