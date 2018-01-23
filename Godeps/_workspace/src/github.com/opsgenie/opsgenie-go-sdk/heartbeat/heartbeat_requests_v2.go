package heartbeat

import (
	"net/url"
	"errors"
)

type PingHeartbeatRequest struct {
	Name   string
	APIKey string
}

func (r *PingHeartbeatRequest) GetApiKey() string {
	return r.APIKey
}

func (r *PingHeartbeatRequest) GenerateUrl() (string, url.Values, error) {
	if r.Name == "" {
		return "", nil, errors.New("Name should be provided")
	}
	return "/v2/heartbeats/" + r.Name + "/ping", nil, nil;
}
