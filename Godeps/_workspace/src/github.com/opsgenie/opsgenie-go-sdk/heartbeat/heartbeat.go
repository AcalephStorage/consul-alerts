package heartbeat

type AddHeartbeatRequest struct {
	ApiKey string 		`json:"apiKey,omitempty"`
	Name string			`json:"name,omitempty"`
	Interval int 			`json:"interval,omitempty"`
	IntervalUnit string	`json:"intervalUnit,omitempty"`
	Description string	`json:"description,omitempty"`
	Enabled bool		`json:"enabled,omitempty"`
}

type UpdateHeartbeatRequest struct {
	ApiKey string		`json:"apiKey,omitempty"`
	Id string 			`json:"id,omitempty"`
	Name string 		`json:"name,omitempty"`
	Interval int 		`json:"interval,omitempty"`
	IntervalUnit string `json:"intervalUnit,omitempty"`
	Description string 	`json:"description,omitempty"`
	Enabled bool 		`json:"enabled,omitempty"`
}

type EnableHeartbeatRequest struct {
	ApiKey string	`url:"apiKey,omitempty"`
	Id string		`url:"id,omitempty"`
	Name string		`url:"name,omitempty"`
}

type DisableHeartbeatRequest struct {
	ApiKey string	`url:"apiKey,omitempty"`
	Id string		`url:"id,omitempty"`
	Name string		`url:"name,omitempty"`
}

type DeleteHeartbeatRequest struct {
	ApiKey string	`url:"apiKey,omitempty"`
	Id string		`url:"id,omitempty"`
	Name string		`url:"name,omitempty"`
}

type GetHeartbeatRequest struct {
	ApiKey string	`url:"apiKey,omitempty"`
	Id string		`url:"id,omitempty"`
	Name string		`url:"name,omitempty"`
}

type ListHeartbeatsRequest struct {
	ApiKey string	`url:"apiKey,omitempty"`
}

type SendHeartbeatRequest struct {
	ApiKey 	string	`json:"apiKey,omitempty"`
	Name 	string	`json:"name,omitempty"` 
}
