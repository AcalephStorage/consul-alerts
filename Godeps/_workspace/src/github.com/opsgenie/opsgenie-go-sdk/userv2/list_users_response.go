package userv2

// ListUsersResponse is a response with list of users.
type ListUsersResponse struct {
	Users []User `json:"data"`
	*ResponseMeta
}
