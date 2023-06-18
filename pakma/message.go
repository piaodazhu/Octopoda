package main

type Response struct {
	Msg         string   `json:"msg"`
	StateType   int      `json:"state_type,omitempty"`
	StateMsg    string   `json:"state_msg,omitempty"`
	Version1    string   `json:"last_version,omitempty"`
	Version2    string   `json:"curr_version,omitempty"`
	Version3    string   `json:"preview_version,omitempty"`
	HistoryList []string `json:"history,omitempty"`
}
