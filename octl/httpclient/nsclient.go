package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/piaodazhu/Octopoda/octl/config"
	"github.com/piaodazhu/Octopoda/octl/output"
	"github.com/piaodazhu/Octopoda/protocols"
	"github.com/piaodazhu/Octopoda/protocols/errs"
)

var nsAddr string

var NsClient *http.Client

func initNsClient() *errs.OctlError {
	NsClient = &http.Client{}
	defaultBrainAddr := fmt.Sprintf("%s:%d", config.GlobalConfig.Brain.Ip, config.GlobalConfig.Brain.Port)
	if !config.GlobalConfig.HttpsNameServer.Enabled {
		output.PrintWarningf("name service client is disabled")
		config.BrainAddr = defaultBrainAddr
		return nil
	}

	config.BrainAddr = ""
	nsAddr = fmt.Sprintf("%s:%d", config.GlobalConfig.HttpsNameServer.Host, config.GlobalConfig.HttpsNameServer.Port+1)
	output.PrintInfof("name server address=%s", nsAddr)

	err := pingNameServer()
	if err != nil {
		emsg := fmt.Sprintf("Could not ping NameServer: %s (%s)", nsAddr, err.Error())
		output.PrintFatalln(emsg)
		return errs.New(errs.OctlInitClientError, emsg)
	}

	entry, err := NameQuery(config.GlobalConfig.Brain.Name + ".octlFace.request")
	if err != nil {
		emsg := fmt.Sprintf("Could not resolve name %s (%s)", config.GlobalConfig.Brain.Name+".octlFace.request", err.Error())
		output.PrintWarningln(emsg)
		return errs.New(errs.OctlInitClientError, emsg)
	}

	config.BrainAddr = entry.Value
	return nil
}

func pingNameServer() error {
	res, err := NsClient.Get(fmt.Sprintf("http://%s/ping", nsAddr))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot Ping https Nameserver")
	}
	return nil
}

func NameQuery(name string) (*protocols.NameServiceEntry, error) {
	res, err := NsClient.Get(fmt.Sprintf("http://%s/query?name=%s", nsAddr, name))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NameQuery status code = %d", res.StatusCode)
	}
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response protocols.Response
	err = json.Unmarshal(buf, &response)
	if err != nil {
		return nil, err
	}
	return response.NameEntry, nil
}
