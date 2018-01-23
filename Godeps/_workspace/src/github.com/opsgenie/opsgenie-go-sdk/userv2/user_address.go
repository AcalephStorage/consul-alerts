package userv2

// UserAddress contains data user user address.
type UserAddress struct {
	Country string `json:"country"`
	State   string `json:"state"`
	City    string `json:"city"`
	Line    string `json:"line"`
	ZipCode string `json:"zipCode"`
}
