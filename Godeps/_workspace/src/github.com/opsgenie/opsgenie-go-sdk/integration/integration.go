package integration

type EnableIntegrationRequest struct {
	ApiKey 	string	`json:"apiKey,omitempty"`
	Id		string 	`json:"id,omitempty"`
	Name 	string	`json:"name,omitempty"`
}

type DisableIntegrationRequest struct {
	ApiKey 	string	`json:"apiKey,omitempty"`
	Id 		string 	`json:"id,omitempty"`
	Name 	string	`json:"name,omitempty"`	
}
