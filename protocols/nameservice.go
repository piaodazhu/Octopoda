package protocols

type NameServiceEntry struct {
	Key         string `json:"key"`
	Type        string `json:"type" binding:"oneof=addr num str array"`
	Value       string `json:"value"`
	Description string `json:"description" binding:"omitempty"`
	TTL         int    `json:"ttl,omitempty" binding:"omitempty"`
}

type ListQueryParam struct {
	Match  string `form:"match" json:"match" binding:"omitempty"`
	Method string `form:"method" json:"method" binding:"oneof=prefix suffix contain equal all"`
}

type Response struct {
	Message   string            `json:"msg,omitempty"`
	NameEntry *NameServiceEntry `json:"entry,omitempty"`
	NameList  []string          `json:"list,omitempty"`
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
