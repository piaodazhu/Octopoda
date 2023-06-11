package main

type NameEntry struct {
	Type        string `redis:"type" json:"type"`
	Name        string `redis:"name" json:"name"`
	Ip          string `redis:"ip" json:"ip"`
	Port        int    `redis:"port" json:"port"`
	TimeStamp   int64  `redis:"ts" json:"ts"`
	Description string `redis:"description" json:"description"`
}

type RegisterParam struct {
	Type        string `form:"type" json:"type" binding:"oneof=brain tentacle other"`
	Name        string `form:"name" json:"name" binding:"required,min=2"`
	Ip          string `form:"ip" json:"ip" binding:"required,ip"`
	Port        int    `form:"port" json:"port" binding:"required,gte=100,lte=65535"`
	Description string `form:"description" json:"description" binding:"required"`
	TTL         int    `form:"ttl" json:"ttl,omitempty" binding:"omitempty"`
}

type ListQueryParam struct {
	Match  string `form:"match" json:"match" binding:"omitempty"`
	Method string `form:"method" json:"method" binding:"oneof=prefix suffix contain equal all"`
}

type Response struct {
	Message   string     `json:"msg,omitempty"`
	NameEntry *NameEntry `json:"entry,omitempty"`
	NameList  []string   `json:"list,omitempty"`
}

type ApiStat struct {
	Requests int `json:"requests"`
	Success  int `json:"success"`
}

type Summary struct {
	TotalRequests int                 `json:"total_request"`
	Since         int64               `json:"since"`
	ApiStats      map[string]*ApiStat `json:"api_stats"`
}