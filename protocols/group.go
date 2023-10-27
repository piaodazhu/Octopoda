package protocols

type GroupInfo struct {
	Name  string   `json:"name" binding:"required"`
	Nodes []string `json:"nodes" binding:"required"`
	// NoCheck can be in request
	NoCheck bool `json:"nocheck" binding:"omitempty"`

	// Size and Unhealthy will be in response
	Size      int      `json:"size" binding:"omitempty"`
	Unhealthy []string `json:"unhealthy" binding:"omitempty"`
}
