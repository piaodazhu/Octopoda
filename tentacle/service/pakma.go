package service

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"tentacle/config"
	"tentacle/logger"
	"tentacle/message"
	"tentacle/snp"
	"unicode"
)

func pakmaState() ([]byte, error) {
	URL := fmt.Sprintf("http://127.0.0.1:%d/state", config.GlobalConfig.PakmaServer.Port)
	res, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	return buf, err
}

func pakmaHistory(timestr string, limit int) ([]byte, error) {
	URL := fmt.Sprintf("http://127.0.0.1:%d/history?time=%s&limit=%d", config.GlobalConfig.PakmaServer.Port, timestr, limit)
	res, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	return buf, err
}

func pakmaInstall(version string) ([]byte, error) {
	URL := fmt.Sprintf("http://127.0.0.1:%d/install", config.GlobalConfig.PakmaServer.Port)
	values := url.Values{}
	values.Set("version", version)
	res, err := http.PostForm(URL, values)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	return buf, err
}

func pakmaUpgrade(version string) ([]byte, error) {
	URL := fmt.Sprintf("http://127.0.0.1:%d/upgrade", config.GlobalConfig.PakmaServer.Port)
	values := url.Values{}
	values.Set("version", version)
	res, err := http.PostForm(URL, values)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	return buf, err
}

func pakmaCancel() ([]byte, error) {
	URL := fmt.Sprintf("http://127.0.0.1:%d/cancel", config.GlobalConfig.PakmaServer.Port)
	res, err := http.PostForm(URL, url.Values{})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	return buf, err
}

func pakmaClean() ([]byte, error) {
	URL := fmt.Sprintf("http://127.0.0.1:%d/clean", config.GlobalConfig.PakmaServer.Port)
	res, err := http.PostForm(URL, url.Values{})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	return buf, err
}

func pakmaConfirm() ([]byte, error) {
	URL := fmt.Sprintf("http://127.0.0.1:%d/confirm", config.GlobalConfig.PakmaServer.Port)
	res, err := http.PostForm(URL, url.Values{})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	return buf, err
}

func pakmaDowngrade() ([]byte, error) {
	URL := fmt.Sprintf("http://127.0.0.1:%d/downgrade", config.GlobalConfig.PakmaServer.Port)
	res, err := http.PostForm(URL, url.Values{})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	buf, err := io.ReadAll(res.Body)
	return buf, err
}

type PakmaParams struct {
	Command string
	Version string
	Time    string
	Limit   int
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

func PakmaCommand(conn net.Conn, raw []byte) {
	rmsg := message.Result{
		Rmsg: "OK",
	}
	var params PakmaParams
	var payload []byte
	var err error
	var valid bool = true
	err = config.Jsoner.Unmarshal(raw, &params)
	if err != nil {
		logger.Exceptions.Println("PakmaCommand Unmarshal")
		rmsg.Rmsg = "PakmaCommand Unmarshal"
		goto errorout
	}

	switch params.Command {
	case "install":
		if !checkVersion(params.Version) {
			valid = false
			logger.Exceptions.Println("PakmaCommand invalid version")
			rmsg.Rmsg = "PakmaCommand invalid version"
			break
		}
		payload, err = pakmaInstall(params.Version)
	case "upgrade":
		if !checkVersion(params.Version) {
			valid = false
			logger.Exceptions.Println("PakmaCommand invalid version")
			rmsg.Rmsg = "PakmaCommand invalid version"
			break
		}
		payload, err = pakmaUpgrade(params.Version)
	case "state":
		payload, err = pakmaState()
	case "confirm":
		payload, err = pakmaConfirm()
	case "cancel":
		payload, err = pakmaCancel()
	case "clean":
		payload, err = pakmaClean()
	case "downgrade":
		payload, err = pakmaDowngrade()
	case "history":
		payload, err = pakmaHistory(params.Time, params.Limit)
	default:
		valid = false
		logger.Exceptions.Println("PakmaCommand unsupport command")
		rmsg.Rmsg = "PakmaCommand unsupport command"
	}
	if !valid {
		payload, _ = config.Jsoner.Marshal(&rmsg)
	} else if err != nil {
		logger.Exceptions.Println("Pakma request error")
		rmsg.Rmsg = "Pakma request error"
		payload, _ = config.Jsoner.Marshal(&rmsg)
	}
errorout:
	err = message.SendMessageUnique(conn, message.TypePakmaCommandResponse, snp.GenSerial(), payload)
	if err != nil {
		logger.Comm.Println("PakmaCommand service error")
	}
}
