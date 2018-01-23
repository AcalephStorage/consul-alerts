/*
Copyright 2015 OpsGenie. All rights reserved.
Use of this source code is governed by a Apache Software
license that can be found in the LICENSE file.
*/

//Package heartbeat provides requests and response structures to achieve Heartbeat API actions.
package heartbeat

import (
	"net/url"
	"errors"
)

// AddHeartbeatRequest provides necessary parameter structure to Create an Heartbeat at OpsGenie.
type AddHeartbeatRequest struct {
	APIKey       string `json:"-"`
	Name         string `json:"name,omitempty"`
	Interval     int    `json:"interval,omitempty"`
	IntervalUnit string `json:"intervalUnit,omitempty"`
	Description  string `json:"description,omitempty"`
	Enabled      *bool  `json:"enabled,omitempty"`
}

func (r *AddHeartbeatRequest) GetApiKey() string {
	return r.APIKey
}

func (r *AddHeartbeatRequest) GenerateUrl() (string, url.Values, error) {
	return "/v2/heartbeats", nil, nil;
}

// UpdateHeartbeatRequest provides necessary parameter structure to Update an existing Heartbeat at OpsGenie.
type UpdateHeartbeatRequest struct {
	APIKey       string `json:"-"`
	Name         string `json:"name,omitempty"`
	Interval     int    `json:"interval,omitempty"`
	IntervalUnit string `json:"intervalUnit,omitempty"`
	Description  string `json:"description,omitempty"`
	Enabled      *bool  `json:"enabled,omitempty"`
}

func (r *UpdateHeartbeatRequest) GetApiKey() string {
	return r.APIKey
}

func (r *UpdateHeartbeatRequest) GenerateUrl() (string, url.Values, error) {
	if r.Name == "" {
		return "", nil, errors.New("Name should be provided")
	}
	return "/v2/heartbeats/" + r.Name, nil, nil;
}

// EnableHeartbeatRequest provides necessary parameter structure to Enable an Heartbeat at OpsGenie.
type EnableHeartbeatRequest struct {
	APIKey string `json:"-"`
	Name   string `json:"name,omitempty"`
}

func (r *EnableHeartbeatRequest) GetApiKey() string {
	return r.APIKey
}

func (r *EnableHeartbeatRequest) GenerateUrl() (string, url.Values, error) {
	if r.Name == "" {
		return "", nil, errors.New("Name should be provided")
	}
	return "/v2/heartbeats/" + r.Name + "/enable", nil, nil;
}

// DisableHeartbeatRequest provides necessary parameter structure to Disable an Heartbeat at OpsGenie.
type DisableHeartbeatRequest struct {
	APIKey string `json:"-"`
	Name   string `json:"name,omitempty"`
}

func (r *DisableHeartbeatRequest) GetApiKey() string {
	return r.APIKey
}

func (r *DisableHeartbeatRequest) GenerateUrl() (string, url.Values, error) {
	if r.Name == "" {
		return "", nil, errors.New("Name should be provided")
	}
	return "/v2/heartbeats/" + r.Name + "/disable", nil, nil;
}

// DeleteHeartbeatRequest provides necessary parameter structure to Delete an Heartbeat from OpsGenie.
type DeleteHeartbeatRequest struct {
	APIKey string `url:"-"`
	Name   string `url:"name,omitempty"`
}

func (r *DeleteHeartbeatRequest) GetApiKey() string {
	return r.APIKey
}

func (r *DeleteHeartbeatRequest) GenerateUrl() (string, url.Values, error) {
	if r.Name == "" {
		return "", nil, errors.New("Name should be provided")
	}
	return "/v2/heartbeats/" + r.Name, nil, nil;
}

// GetHeartbeatRequest provides necessary parameter structure to Retrieve an Heartbeat with details from OpsGenie.
type GetHeartbeatRequest struct {
	APIKey string `url:"-"`
	Name   string `url:"name,omitempty"`
}


func (r *GetHeartbeatRequest) GetApiKey() string {
	return r.APIKey
}

func (r *GetHeartbeatRequest) GenerateUrl() (string, url.Values, error) {
	if r.Name == "" {
		return "", nil, errors.New("Name should be provided")
	}
	return "/v2/heartbeats/" + r.Name, nil, nil;
}

// ListHeartbeatsRequest provides necessary parameter structure to Retrieve Heartbeats from OpsGenie.
type ListHeartbeatsRequest struct {
	APIKey string `url:"apiKey"`
}

// SendHeartbeatRequest provides necessary parameter structure to Send an Heartbeat Signal to OpsGenie.
type SendHeartbeatRequest struct {
	APIKey string `json:"apiKey"`
	Name   string `json:"name,omitempty"`
}
