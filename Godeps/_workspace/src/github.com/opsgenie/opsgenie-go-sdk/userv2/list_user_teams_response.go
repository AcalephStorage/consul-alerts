package userv2

// ListUserTeamsResponse is a response with list of user teams.
type ListUserTeamsResponse struct {
	Teams []Team `json:"data"`
	ResponseMeta
}
