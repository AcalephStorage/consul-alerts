package policy

type EnablePolicyResponse struct {
	Status string 	`json:"status"`
	Code int 		`json:"code"`
}

type DisablePolicyResponse struct {
	Status string 	`json:"status"`
	Code int  		`json:"code"`
}
