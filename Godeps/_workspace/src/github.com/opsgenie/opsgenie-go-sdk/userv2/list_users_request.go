package userv2

import (
	"errors"
	"net/url"
	"strconv"
)

// ListUsersRequest is a request for getting user list.
type ListUsersRequest struct {
	Limit  int
	Offset int
	Sort   Sort
	Order  Order
	Query  string
	ApiKey string
}

// GenerateUrl generates API url for getting user list.
func (r *ListUsersRequest) GenerateUrl() (string, url.Values, error) {
	baseUrl := "/v2/users/"

	params := url.Values{}
	params.Set("offset", strconv.Itoa(r.Offset))

	if r.Limit > 0 {
		params.Set("limit", strconv.Itoa(r.Limit))
	}

	if len(r.Sort) > 0 {
		if r.Sort.IsValid() {
			params.Set("sort", string(r.Sort))
		} else {
			return "", nil, errors.New("unavailable field to use in sorting, id should be one of 'username', 'fullName' or 'insertedAt'")
		}
	}

	if len(r.Order) > 0 {
		if r.Order.IsValid() {
			params.Set("order", string(r.Order))
		} else {
			return "", nil, errors.New("unavailable direction of sorting, it should be 'asc' or 'desc'")
		}
	}

	if len(r.Query) > 0 {
		params.Set("query", r.Query)
	}

	return baseUrl, params, nil
}

// GetApiKey returns api key.
func (r *ListUsersRequest) GetApiKey() string {
	return r.ApiKey
}
