package userv2

// GetUserResponse is a response of getting user result.
type GetUserResponse struct {
	User       User     `json:"data"`
	ResponseMeta
}
