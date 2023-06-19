package node

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"octl/config"
	"octl/nameclient"
	"octl/output"
	"unicode"
)

func Pakma(firstarg string, args []string) {
	var subcmd string = firstarg
	var version string = ""
	var nodes []string

	if len(args) < 1 {
		output.PrintFatalln("not any node is specified")
		return
	}

	switch firstarg {
	case "install":
		fallthrough // same process as upgrade
	case "upgrade":
		if len(args) < 2 {
			output.PrintFatalln("not any node is specified")
			return
		}
		version = args[0]
		nodes = args[1:]
	case "state":
		fallthrough
	case "confirm":
		fallthrough
	case "cancel":
		fallthrough
	case "downgrade":
		nodes = args
	default:
		output.PrintFatalln("pakma subcommand not support: ", firstarg)
		return
	}

	url := fmt.Sprintf("http://%s/%s%s",
		nameclient.BrainAddr,
		config.GlobalConfig.Brain.ApiPrefix,
		config.GlobalConfig.Api.Pakma,
	)
	nodes_serialized, _ := config.Jsoner.Marshal(&nodes)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()
	writer.WriteField("subcmd", subcmd)
	writer.WriteField("targetNodes", string(nodes_serialized))

	if version != "" {
		if checkVersion(version) {
			writer.WriteField("version", version)
		} else {
			output.PrintFatalln("pakma invalid version number (right example: 1.4.1): ", version)
			return
		}
	}

	res, err := http.Post(url, writer.FormDataContentType(), body)
	if err != nil {
		output.PrintFatalln("Post")
	}
	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)
	output.PrintJSON(raw)
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
