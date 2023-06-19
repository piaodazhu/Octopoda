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
	"time"
)

func pakmaState() ([]byte, error) {
	URL := fmt.Sprintf("http://127.0.0.1:%d/state", config.GlobalConfig.PakmaServer.Port)
	res, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	res.Body.Close()
	return io.ReadAll(res.Body)
}

func pakmaHistory(timestr string, limit int) ([]byte, error) {
	URL := fmt.Sprintf("http://127.0.0.1:%d/history?time=%s&limit=%d", config.GlobalConfig.PakmaServer.Port, timestr, limit)
	res, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	res.Body.Close()
	return io.ReadAll(res.Body)
}

func pakmaInstall(version string) ([]byte, error) {
	URL := fmt.Sprintf("http://127.0.0.1:%d/install", config.GlobalConfig.PakmaServer.Port)
	values := url.Values{}
	values.Set("version", version)
	res, err := http.PostForm(URL, values)
	if err != nil {
		return nil, err
	}
	res.Body.Close()
	return io.ReadAll(res.Body)
}

func pakmaUpgrade(version string) ([]byte, error) {
	URL := fmt.Sprintf("http://127.0.0.1:%d/upgrade", config.GlobalConfig.PakmaServer.Port)
	values := url.Values{}
	values.Set("version", version)
	res, err := http.PostForm(URL, values)
	if err != nil {
		return nil, err
	}
	res.Body.Close()
	return io.ReadAll(res.Body)
}

func pakmaCancel() ([]byte, error) {
	URL := fmt.Sprintf("http://127.0.0.1:%d/cancel", config.GlobalConfig.PakmaServer.Port)
	res, err := http.PostForm(URL, url.Values{})
	if err != nil {
		return nil, err
	}
	res.Body.Close()
	return io.ReadAll(res.Body)
}

func pakmaConfirm() ([]byte, error) {
	URL := fmt.Sprintf("http://127.0.0.1:%d/confirm", config.GlobalConfig.PakmaServer.Port)
	res, err := http.PostForm(URL, url.Values{})
	if err != nil {
		return nil, err
	}
	res.Body.Close()
	return io.ReadAll(res.Body)
}

func pakmaDowngrade() ([]byte, error) {
	URL := fmt.Sprintf("http://127.0.0.1:%d/downgrade", config.GlobalConfig.PakmaServer.Port)
	res, err := http.PostForm(URL, url.Values{})
	if err != nil {
		return nil, err
	}
	res.Body.Close()
	return io.ReadAll(res.Body)
}

func PakmaCommand(conn net.Conn, raw []byte) {
	stateLock.RLock()
	state := nodeStatus
	state.LocalTime = time.Now().UnixNano()
	stateLock.RUnlock()
	serialized_info, _ := config.Jsoner.Marshal(&state)
	err := message.SendMessage(conn, message.TypeNodeStatusResponse, serialized_info)
	if err != nil {
		logger.Comm.Println("NodeStatus service error")
	}
}
