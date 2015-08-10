package integration

type EnableIntegrationResponse struct {
	Status string 	`json:"status"`
	Code int 		`json:"code"`
}

type DisableIntegrationResponse struct {
	Status string 	`json:"status"`
	Code int  		`json:"code"`
}
