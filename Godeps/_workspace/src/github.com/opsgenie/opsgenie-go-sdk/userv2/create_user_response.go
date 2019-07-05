package userv2

// CreateUserResponse is a response of creating user result.
type CreateUserResponse struct {
	User UserMeta `json:"data"`
	ResponseMeta
}
