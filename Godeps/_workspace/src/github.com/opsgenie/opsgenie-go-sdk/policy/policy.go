package policy

type EnablePolicyRequest struct {
	ApiKey 	string	`json:"apiKey,omitempty"`
	Id		string 	`json:"id,omitempty"`
	Name 	string	`json:"name,omitempty"`
}

type DisablePolicyRequest struct {
	ApiKey 	string	`json:"apiKey,omitempty"`
	Id 		string 	`json:"id,omitempty"`
	Name 	string	`json:"name,omitempty"`	
}
