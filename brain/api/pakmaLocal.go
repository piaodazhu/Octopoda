package api

import (
	"brain/config"
	"fmt"
	"io"
	"net/http"
	"net/url"

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
