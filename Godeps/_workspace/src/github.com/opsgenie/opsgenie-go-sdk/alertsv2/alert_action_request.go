package alertsv2

import (
	"net/url"
	"time"
)

type AlertActionRequest struct {
	*Identifier
	User   string `json:"user,omitempty"`
	Source string `json:"source,omitempty"`
	Note   string `json:"note,omitempty"`
	ApiKey string `json:"-"`
}

func (r *AlertActionRequest) GetApiKey() string {
	return r.ApiKey
}

type AcknowledgeRequest AlertActionRequest

func (r *AcknowledgeRequest) GenerateUrl() (string, url.Values, error) {
	path, params, err := r.Identifier.GenerateUrl()
	return path + "/acknowledge", params, err
}

func (r *AcknowledgeRequest) GetApiKey() string {
	return r.ApiKey
}

type CloseRequest AlertActionRequest

func (r *CloseRequest) GenerateUrl() (string, url.Values, error) {
	path, params, err := r.Identifier.GenerateUrl()
	return path + "/close", params, err
}

func (r *CloseRequest) GetApiKey() string {
	return r.ApiKey
}

type UnacknowledgeRequest AlertActionRequest

func (r *UnacknowledgeRequest) GenerateUrl() (string, url.Values, error) {
	path, params, err := r.Identifier.GenerateUrl()
	return path + "/unacknowledge", params, err
}

func (r *UnacknowledgeRequest) GetApiKey() string {
	return r.ApiKey
}

type SnoozeRequest struct {
	AlertActionRequest
	EndTime time.Time `json:"endTime"`
}

func (r *SnoozeRequest) GenerateUrl() (string, url.Values, error) {
	path, params, err := r.Identifier.GenerateUrl()
	return path + "/snooze", params, err
}

func (r *SnoozeRequest) GetApiKey() string {
	return r.ApiKey
}

type AddNoteRequest AlertActionRequest

func (r *AddNoteRequest) GenerateUrl() (string, url.Values, error) {
	path, params, err := r.Identifier.GenerateUrl()
	return path + "/notes", params, err
}

func (r *AddNoteRequest) GetApiKey() string {
	return r.ApiKey
}
