package userv2

// ListUserEscalationsResponse is a response with list of user escalations.
type ListUserEscalationsResponse struct {
	Escalations []Escalation `json:"data,omitempty"`
	ResponseMeta
}
