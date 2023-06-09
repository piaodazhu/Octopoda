package nameclient

type NameEntry struct {
	RegisterParam
	TimeStamp int64 `json:"ts"`
}

type RegisterParam struct {
	Type        string `form:"type" json:"type" binding:"oneof=brain tentacle octl other"`
	Name        string `form:"name" json:"name" binding:"required,min=2"`
	Ip          string `form:"ip" json:"ip" binding:"required,ip"`
	Port        int    `form:"port" json:"port" binding:"required,gte=1,lte=65535"`
	Port2       int    `form:"port2" json:"port2" binding:"omitemtpy,gte=1,lte=65535"`
	Description string `form:"description" json:"description" binding:"omitempty"`
	TTL         int    `form:"ttl" json:"ttl,omitempty" binding:"omitempty"`
}

type ListQueryParam struct {
	Scope  string `form:"scope" json:"scope" binding:"oneof=name config ssh"`
	Match  string `form:"match" json:"match" binding:"omitempty"`
	Method string `form:"method" json:"method" binding:"oneof=prefix suffix contain equal all"`
}

type ConfigEntry struct {
	ConfigUploadParam
	TimeStamp int64 `json:"ts"`
}

type ConfigUploadParam struct {
	Type      string `form:"type" json:"type" binding:"oneof=brain tentacle octl other"`
	Name      string `form:"name" json:"name" binding:"required,min=2"`
	Method    string `form:"method" json:"method,omitempty" binding:"oneof=reset append clear"`
	RawConfig string `form:"conf" json:"conf" binding:"omitempty"`
}

type ConfigQueryParam struct {
	Name   string `form:"name" json:"name" binding:"required,min=2"`
	Index  int    `form:"index" json:"index" binding:"gte=0"`
	Amount int    `form:"amount" json:"amount" binding:"required,gte=1"`
}

type SshInfo struct {
	SshInfoUploadParam
	TimeStamp int64 `json:"ts"`
}

type SshInfoUploadParam struct {
	Type     string `form:"type" json:"type" binding:"oneof=brain tentacle octl other"`
	Name     string `form:"name" json:"name" binding:"required,min=2"`
	Username string `form:"username" json:"username" binding:"required"`
	Ip       string `form:"ip" json:"ip" binding:"required,ip"`
	Port     int    `form:"port" json:"port" binding:"required,gte=1,lte=65535"`
	Password string `form:"password" json:"password" binding:"required"`
}

type Response struct {
	Message   string         `json:"msg,omitempty"`
	NameEntry *NameEntry     `json:"entry,omitempty"`
	NameList  []string       `json:"list,omitempty"`
	RawConfig []*ConfigEntry `json:"conflist,omitempty"`
	SshInfo   *SshInfo       `json:"sshinfo,omitempty"`
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
