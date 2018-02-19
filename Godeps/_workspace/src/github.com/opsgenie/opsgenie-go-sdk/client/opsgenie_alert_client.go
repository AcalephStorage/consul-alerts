// Copyright 2015 OpsGenie. All rights reserved.
// Use of this source code is governed by an Apache Software
// license that can be found in the LICENSE file.
package client

import (
	"bytes"
	"errors"
	"fmt"
	goreq "github.com/Difrex/consul-alerts/Godeps/_workspace/src/github.com/franela/goreq"
	goquery "github.com/Difrex/consul-alerts/Godeps/_workspace/src/github.com/google/go-querystring/query"
	alerts "github.com/Difrex/consul-alerts/Godeps/_workspace/src/github.com/opsgenie/opsgenie-go-sdk/alerts"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	CREATE_ALERT_URL           = ENDPOINT_URL + "/v1/json/alert"
	CLOSE_ALERT_URL            = ENDPOINT_URL + "/v1/json/alert/close"
	DELETE_ALERT_URL           = ENDPOINT_URL + "/v1/json/alert"
	GET_ALERT_URL              = ENDPOINT_URL + "/v1/json/alert"
	LIST_ALERTS_URL            = ENDPOINT_URL + "/v1/json/alert"
	LIST_ALERT_NOTES_URL       = ENDPOINT_URL + "/v1/json/alert/note"
	LIST_ALERT_LOGS_URL        = ENDPOINT_URL + "/v1/json/alert/log"
	LIST_ALERT_RECIPIENTS_URL  = ENDPOINT_URL + "/v1/json/alert/recipient"
	ACKNOWLEDGE_ALERT_URL      = ENDPOINT_URL + "/v1/json/alert/acknowledge"
	RENOTIFY_ALERT_URL         = ENDPOINT_URL + "/v1/json/alert/renotify"
	TAKE_OWNERSHIP_ALERT_URL   = ENDPOINT_URL + "/v1/json/alert/takeOwnership"
	ASSIGN_OWNERSHIP_ALERT_URL = ENDPOINT_URL + "/v1/json/alert/assign"
	ADD_TEAM_ALERT_URL         = ENDPOINT_URL + "/v1/json/alert/team"
	ADD_RECIPIENT_ALERT_URL    = ENDPOINT_URL + "/v1/json/alert/recipient"
	ADD_NOTE_ALERT_URL         = ENDPOINT_URL + "/v1/json/alert/note"
	EXECUTE_ACTION_ALERT_URL   = ENDPOINT_URL + "/v1/json/alert/executeAction"
	ATTACH_FILE_ALERT_URL      = ENDPOINT_URL + "/v1/json/alert/attach"
)

type OpsGenieAlertClient struct {
	apiKey  string
	proxy   string
	retries int
}

func (cli *OpsGenieAlertClient) buildRequest(method string, uri string, body interface{}) goreq.Request {
	req := goreq.Request{}
	req.Method = method
	req.Uri = uri
	if body != nil {
		req.Body = body
	}
	if cli.proxy != "" {
		req.Proxy = cli.proxy
	}
	req.UserAgent = userAgentParam.ToString()
	return req
}

func (cli *OpsGenieAlertClient) SetConnectionTimeout(timeoutInSeconds time.Duration) {
	goreq.SetConnectTimeout(timeoutInSeconds * time.Second)
}

func (cli *OpsGenieAlertClient) SetMaxRetryAttempts(retries int) {
	cli.retries = retries
}

func (cli *OpsGenieAlertClient) Create(req alerts.CreateAlertRequest) (*alerts.CreateAlertResponse, error) {
	req.ApiKey = cli.apiKey
	// validate the mandatory parameters: apiKey, message
	if req.ApiKey == "" {
		return nil, errors.New("ApiKey is a mandatory field and can not be empty.")
	}
	if req.Message == "" {
		return nil, errors.New("Message is a mandatory field and can not be empty.")
	}
	// send the request
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("POST", CREATE_ALERT_URL, req).Do()
		if err == nil {
			break
		}
		// sleep for a second
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Could not create the alert: a problem occured while sending the request.")
	}
	// check the returning HTTP status code
	httpStatusCode := resp.StatusCode
	if httpStatusCode >= 400 && httpStatusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d returned.", httpStatusCode))
	}
	if httpStatusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d returned.", httpStatusCode))
	}
	var createAlertResp alerts.CreateAlertResponse
	// check if the response can be unmarshalled
	if err = resp.Body.FromJsonTo(&createAlertResp); err != nil {
		return nil, errors.New("Server response can not be parsed")
	}
	return &createAlertResp, nil
}

func (cli *OpsGenieAlertClient) Close(req alerts.CloseAlertRequest) (*alerts.CloseAlertResponse, error) {
	req.ApiKey = cli.apiKey
	// validate the mandatory parameters: apiKey, alertId/alias
	if req.ApiKey == "" {
		return nil, errors.New("ApiKey is a mandatory field and can not be empty.")
	}
	if req.AlertId == "" && req.Alias == "" {
		return nil, errors.New("At least either Alert Id or Alias should be set in the request.")
	}
	if req.AlertId != "" && req.Alias != "" {
		return nil, errors.New("Either Alert Id or Alias should be set in the request not both.")
	}
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("POST", CLOSE_ALERT_URL, req).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	// resp, err := goreq.Request{ Method: "POST", Uri: CLOSE_ALERT_URL, Body: req, }.Do()
	if err != nil {
		return nil, errors.New("Could not close the alert: a problem occured while sending the request.")
	}
	// check the returning HTTP status code
	httpStatusCode := resp.StatusCode
	if httpStatusCode >= 400 && httpStatusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d returned", httpStatusCode))
	}
	if httpStatusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d returned", httpStatusCode))
	}
	var closeAlertResp alerts.CloseAlertResponse
	// check if the response can be unmarshalled
	if err = resp.Body.FromJsonTo(&closeAlertResp); err != nil {
		return nil, errors.New("Server response can not be parsed")
	}
	return &closeAlertResp, nil
}

func (cli *OpsGenieAlertClient) Delete(req alerts.DeleteAlertRequest) (*alerts.DeleteAlertResponse, error) {
	req.ApiKey = cli.apiKey
	// validate the mandatory parameters: apiKey, alertId/alias
	if req.ApiKey == "" {
		return nil, errors.New("ApiKey is a mandatory field and can not be empty.")
	}
	if req.AlertId == "" && req.Alias == "" {
		return nil, errors.New("Either Alert Id or Alias at least should be set in the request.")
	}
	if req.AlertId != "" && req.Alias != "" {
		return nil, errors.New("Either Alert Id or Alias should be set in the request not both.")
	}
	v, _ := goquery.Values(req)
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("DELETE", DELETE_ALERT_URL+"?"+v.Encode(), nil).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Could not delete the alert: a problem occured while sending the request.")
	}
	// check the returning HTTP status code
	httpStatusCode := resp.StatusCode
	if httpStatusCode >= 400 && httpStatusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d returned", httpStatusCode))
	}
	if httpStatusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d returned", httpStatusCode))
	}
	var deleteAlertResp alerts.DeleteAlertResponse
	// check if the response can be unmarshalled
	if err = resp.Body.FromJsonTo(&deleteAlertResp); err != nil {
		return nil, errors.New("Server response can not be parsed.")
	}
	return &deleteAlertResp, nil
}

func (cli *OpsGenieAlertClient) Get(req alerts.GetAlertRequest) (*alerts.GetAlertResponse, error) {
	req.ApiKey = cli.apiKey
	// validate the mandatory parameters: apiKey, id/alias/tinyId
	if req.ApiKey == "" {
		return nil, errors.New("ApiKey is a mandatory field and can not be empty.")
	}
	if req.Id == "" && req.Alias == "" && req.TinyId == "" {
		return nil, errors.New("At least one of the parameters of id, alias and tiny id should be set.")
	}
	if (req.Id != "" && req.Alias != "") || (req.Id != "" && req.TinyId != "") || (req.Alias != "" && req.TinyId != "") {
		return nil, errors.New("Only one of the parameters of id, alias and tiny id should be set.")
	}
	v, _ := goquery.Values(req)
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("GET", GET_ALERT_URL+"?"+v.Encode(), nil).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Could not retrieve the alert: a problem occured while sending the request")
	}
	// check the returning HTTP status code
	httpStatusCode := resp.StatusCode
	if httpStatusCode >= 400 && httpStatusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d returned", httpStatusCode))
	}
	if httpStatusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d returned", httpStatusCode))
	}
	var getAlertResp alerts.GetAlertResponse
	// check if the response can be unmarshalled
	if err = resp.Body.FromJsonTo(&getAlertResp); err != nil {
		return nil, errors.New("Server response can not be parsed.")
	}
	return &getAlertResp, nil
}

func (cli *OpsGenieAlertClient) List(req alerts.ListAlertsRequest) (*alerts.ListAlertsResponse, error) {
	req.ApiKey = cli.apiKey
	// validate the mandatory parameter: apiKey
	if req.ApiKey == "" {
		return nil, errors.New("ApiKey is a mandatory field and can not be empty.")
	}
	v, _ := goquery.Values(req)
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("GET", LIST_ALERTS_URL+"?"+v.Encode(), nil).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Could not retrieve the alert: a problem occured while sending the request")
	}
	// check the returning HTTP status code
	httpStatusCode := resp.StatusCode
	if httpStatusCode >= 400 && httpStatusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d returned", httpStatusCode))
	}
	if httpStatusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d returned", httpStatusCode))
	}
	var listAlertsResp alerts.ListAlertsResponse
	// check if the response can be unmarshalled
	if err = resp.Body.FromJsonTo(&listAlertsResp); err != nil {
		return nil, errors.New("Server response can not be parsed.")
	}
	return &listAlertsResp, nil
}

func (cli *OpsGenieAlertClient) ListNotes(req alerts.ListAlertNotesRequest) (*alerts.ListAlertNotesResponse, error) {
	req.ApiKey = cli.apiKey
	// validate the mandatory parameters: apiKey, id/alias
	if req.ApiKey == "" {
		return nil, errors.New("ApiKey is a mandatory field and can not be empty.")
	}
	if req.Id == "" && req.Alias == "" {
		return nil, errors.New("At least either Id or Alias should be set in the request.")
	}
	if req.Id != "" && req.Alias != "" {
		return nil, errors.New("Either Id or Alias should be set in the request not both.")
	}
	v, _ := goquery.Values(req)
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("GET", LIST_ALERT_NOTES_URL+"?"+v.Encode(), nil).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Could not send the request: a problem occured while sending the request")
	}
	// check the returning HTTP status code
	httpStatusCode := resp.StatusCode
	if httpStatusCode >= 400 && httpStatusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d returned", httpStatusCode))
	}
	if httpStatusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d returned", httpStatusCode))
	}
	var listAlertNotesResp alerts.ListAlertNotesResponse
	// check if the response can be unmarshalled
	if err = resp.Body.FromJsonTo(&listAlertNotesResp); err != nil {
		return nil, errors.New("Server response can not be parsed.")
	}
	return &listAlertNotesResp, nil
}

func (cli *OpsGenieAlertClient) ListLogs(req alerts.ListAlertLogsRequest) (*alerts.ListAlertLogsResponse, error) {
	req.ApiKey = cli.apiKey
	// validate the mandatory parameters: apiKey, id/alias
	if req.ApiKey == "" {
		return nil, errors.New("ApiKey is a mandatory field and can not be empty.")
	}
	if req.Id == "" && req.Alias == "" {
		return nil, errors.New("At least either Id or Alias should be set in the request.")
	}
	if req.Id != "" && req.Alias != "" {
		return nil, errors.New("Either Id or Alias should be set in the request not both.")
	}
	v, _ := goquery.Values(req)
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("GET", LIST_ALERT_LOGS_URL+"?"+v.Encode(), nil).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Could not retrieve the logs: a problem occured while sending the request")
	}
	// check the returning HTTP status code
	httpStatusCode := resp.StatusCode
	if httpStatusCode >= 400 && httpStatusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d returned", httpStatusCode))
	}
	if httpStatusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d returned", httpStatusCode))
	}
	var listAlertLogsResp alerts.ListAlertLogsResponse
	// check if the response can be unmarshalled
	if err = resp.Body.FromJsonTo(&listAlertLogsResp); err != nil {
		return nil, errors.New("Server response can not be parsed.")
	}
	return &listAlertLogsResp, nil
}

func (cli *OpsGenieAlertClient) ListRecipients(req alerts.ListAlertRecipientsRequest) (*alerts.ListAlertRecipientsResponse, error) {
	req.ApiKey = cli.apiKey
	// validate the mandatory parameters: apiKey, id/alias
	if req.ApiKey == "" {
		return nil, errors.New("ApiKey is a mandatory field and can not be empty.")
	}
	if req.Id == "" && req.Alias == "" {
		return nil, errors.New("At least either Id or Alias should be set in the request.")
	}
	if req.Id != "" && req.Alias != "" {
		return nil, errors.New("Either Id or Alias should be set in the request not both.")
	}
	v, _ := goquery.Values(req)
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("GET", LIST_ALERT_RECIPIENTS_URL+"?"+v.Encode(), nil).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Can not list the recipient list, unable to send the request")
	}
	// check the returning HTTP status code
	httpStatusCode := resp.StatusCode
	if httpStatusCode >= 400 && httpStatusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d returned", httpStatusCode))
	}
	if httpStatusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d returned", httpStatusCode))
	}
	var listAlertRecipientsResp alerts.ListAlertRecipientsResponse
	// check if the response can be unmarshalled
	if err = resp.Body.FromJsonTo(&listAlertRecipientsResp); err != nil {
		return nil, errors.New("Server response can not be parsed.")
	}
	return &listAlertRecipientsResp, nil
}

func (cli *OpsGenieAlertClient) Acknowledge(req alerts.AcknowledgeAlertRequest) (*alerts.AcknowledgeAlertResponse, error) {
	req.ApiKey = cli.apiKey
	// validate the mandatory parameters: apiKey, alertId/alias
	if req.ApiKey == "" {
		return nil, errors.New("ApiKey is a mandatory field and can not be empty.")
	}
	if req.AlertId == "" && req.Alias == "" {
		return nil, errors.New("At least either Alert Id or Alias should be set in the request.")
	}
	if req.AlertId != "" && req.Alias != "" {
		return nil, errors.New("Either Alert Id or Alias should be set in the request not both.")
	}
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("POST", ACKNOWLEDGE_ALERT_URL, req).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Can not ack the alert, unable to send the request")
	}
	// check the returning HTTP status code
	httpStatusCode := resp.StatusCode
	if httpStatusCode >= 400 && httpStatusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d returned", httpStatusCode))
	}
	if httpStatusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d returned", httpStatusCode))
	}
	var acknowledgeAlertResp alerts.AcknowledgeAlertResponse
	// check if the response can be unmarshalled
	if err = resp.Body.FromJsonTo(&acknowledgeAlertResp); err != nil {
		return nil, errors.New("Server response can not be parsed.")
	}
	return &acknowledgeAlertResp, nil
}

func (cli *OpsGenieAlertClient) Renotify(req alerts.RenotifyAlertRequest) (*alerts.RenotifyAlertResponse, error) {
	req.ApiKey = cli.apiKey
	// validate the mandatory parameters: apiKey, alertId/alias
	if req.ApiKey == "" {
		return nil, errors.New("ApiKey is a mandatory field and can not be empty.")
	}
	if req.AlertId == "" && req.Alias == "" {
		return nil, errors.New("At least either Alert Id or Alias should be set in the request.")
	}
	if req.AlertId != "" && req.Alias != "" {
		return nil, errors.New("Either Alert Id or Alias should be set in the request not both.")
	}
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("POST", RENOTIFY_ALERT_URL, req).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Can not renotify, unable to send the request")
	}
	// check the returning HTTP status code
	httpStatusCode := resp.StatusCode
	if httpStatusCode >= 400 && httpStatusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d returned", httpStatusCode))
	}
	if httpStatusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d returned", httpStatusCode))
	}
	var renotifyAlertResp alerts.RenotifyAlertResponse
	// check if the response can be unmarshalled
	if err = resp.Body.FromJsonTo(&renotifyAlertResp); err != nil {
		return nil, errors.New("Server response can not be parsed.")
	}
	return &renotifyAlertResp, nil
}

func (cli *OpsGenieAlertClient) TakeOwnership(req alerts.TakeOwnershipAlertRequest) (*alerts.TakeOwnershipAlertResponse, error) {
	req.ApiKey = cli.apiKey
	// validate the mandatory parameters: apiKey, alertId/alias
	if req.ApiKey == "" {
		return nil, errors.New("ApiKey is a mandatory field and can not be empty.")
	}
	if req.AlertId == "" && req.Alias == "" {
		return nil, errors.New("At least either Alert Id or Alias should be set in the request.")
	}
	if req.AlertId != "" && req.Alias != "" {
		return nil, errors.New("Either Alert Id or Alias should be set in the request not both.")
	}
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("POST", TAKE_OWNERSHIP_ALERT_URL, req).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Can not change the ownership, unable to send the request")
	}
	// check the returning HTTP status code
	httpStatusCode := resp.StatusCode
	if httpStatusCode >= 400 && httpStatusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d returned", httpStatusCode))
	}
	if httpStatusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d returned", httpStatusCode))
	}
	var takeOwnershipResp alerts.TakeOwnershipAlertResponse
	// check if the response can be unmarshalled
	if err = resp.Body.FromJsonTo(&takeOwnershipResp); err != nil {
		return nil, errors.New("Server response can not be parsed.")
	}
	return &takeOwnershipResp, nil
}

func (cli *OpsGenieAlertClient) AssignOwner(req alerts.AssignOwnerAlertRequest) (*alerts.AssignOwnerAlertResponse, error) {
	req.ApiKey = cli.apiKey
	// validate the mandatory parameters: apiKey, alertId/alias, owner
	if req.ApiKey == "" {
		return nil, errors.New("ApiKey is a mandatory field and can not be empty.")
	}
	if req.Owner == "" {
		return nil, errors.New("Owner is a mandatory field and can not be empty")
	}
	if req.AlertId == "" && req.Alias == "" {
		return nil, errors.New("At least either Alert Id or Alias should be set in the request.")
	}
	if req.AlertId != "" && req.Alias != "" {
		return nil, errors.New("Either Alert Id or Alias should be set in the request not both.")
	}
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("POST", ASSIGN_OWNERSHIP_ALERT_URL, req).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Can not assign the owner, unable to send the request")
	}
	// check the returning HTTP status code
	httpStatusCode := resp.StatusCode
	if httpStatusCode >= 400 && httpStatusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d returned", httpStatusCode))
	}
	if httpStatusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d returned", httpStatusCode))
	}
	var assignOwnerAlertResp alerts.AssignOwnerAlertResponse
	// check if the response can be unmarshalled
	if err = resp.Body.FromJsonTo(&assignOwnerAlertResp); err != nil {
		return nil, errors.New("Server response can not be parsed.")
	}
	return &assignOwnerAlertResp, nil
}

func (cli *OpsGenieAlertClient) AddTeam(req alerts.AddTeamAlertRequest) (*alerts.AddTeamAlertResponse, error) {
	req.ApiKey = cli.apiKey
	// validate the mandatory parameters: apiKey, alertId/alias, team
	if req.ApiKey == "" {
		return nil, errors.New("ApiKey is a mandatory field and can not be empty.")
	}
	if req.Team == "" {
		return nil, errors.New("Team is a mandatory field and can not be empty")
	}
	if req.AlertId == "" && req.Alias == "" {
		return nil, errors.New("At least either Alert Id or Alias should be set in the request.")
	}
	if req.AlertId != "" && req.Alias != "" {
		return nil, errors.New("Either Alert Id or Alias should be set in the request not both.")
	}
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("POST", ADD_TEAM_ALERT_URL, req).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Team can not be added, unable to send the request")
	}
	// check the returning HTTP status code
	httpStatusCode := resp.StatusCode
	if httpStatusCode >= 400 && httpStatusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d returned", httpStatusCode))
	}
	if httpStatusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d returned", httpStatusCode))
	}
	var addTeamAlertResp alerts.AddTeamAlertResponse
	// check if the response can be unmarshalled
	if err = resp.Body.FromJsonTo(&addTeamAlertResp); err != nil {
		return nil, errors.New("Server response can not be parsed.")
	}
	return &addTeamAlertResp, nil
}

func (cli *OpsGenieAlertClient) AddRecipient(req alerts.AddRecipientAlertRequest) (*alerts.AddRecipientAlertResponse, error) {
	req.ApiKey = cli.apiKey
	// validate the mandatory parameters: apiKey, alertId/alias, recipient
	if req.ApiKey == "" {
		return nil, errors.New("ApiKey is a mandatory field and can not be empty.")
	}
	if req.Recipient == "" {
		return nil, errors.New("Recipient is a mandatory field and can not be empty")
	}
	if req.AlertId == "" && req.Alias == "" {
		return nil, errors.New("At least either Alert Id or Alias should be set in the request.")
	}
	if req.AlertId != "" && req.Alias != "" {
		return nil, errors.New("Either Alert Id or Alias should be set in the request not both.")
	}
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("POST", ADD_RECIPIENT_ALERT_URL, req).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Can not add recipient, unable to send the request")
	}
	// check the returning HTTP status code
	httpStatusCode := resp.StatusCode
	if httpStatusCode >= 400 && httpStatusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d returned", httpStatusCode))
	}
	if httpStatusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d returned", httpStatusCode))
	}
	var addRecipientAlertResp alerts.AddRecipientAlertResponse
	// check if the response can be unmarshalled
	if err = resp.Body.FromJsonTo(&addRecipientAlertResp); err != nil {
		return nil, errors.New("Server response can not be parsed.")
	}
	return &addRecipientAlertResp, nil
}

func (cli *OpsGenieAlertClient) AddNote(req alerts.AddNoteAlertRequest) (*alerts.AddNoteAlertResponse, error) {
	req.ApiKey = cli.apiKey
	// validate the mandatory parameters: apiKey, alertId/alias, note
	if req.ApiKey == "" {
		return nil, errors.New("ApiKey is a mandatory field and can not be empty.")
	}
	if req.Note == "" {
		return nil, errors.New("Note is a mandatory field and can not be empty")
	}
	if req.AlertId == "" && req.Alias == "" {
		return nil, errors.New("At least either Alert Id or Alias should be set in the request.")
	}
	if req.AlertId != "" && req.Alias != "" {
		return nil, errors.New("Either Alert Id or Alias should be set in the request not both.")
	}
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("POST", ADD_NOTE_ALERT_URL, req).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Can not add note, unable to send the request.")
	}
	// check the returning HTTP status code
	httpStatusCode := resp.StatusCode
	if httpStatusCode >= 400 && httpStatusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d returned", httpStatusCode))
	}
	if httpStatusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d returned", httpStatusCode))
	}
	var addNoteAlertResp alerts.AddNoteAlertResponse
	// check if the response can be unmarshalled
	if err = resp.Body.FromJsonTo(&addNoteAlertResp); err != nil {
		return nil, errors.New("Server response can not be parsed.")
	}
	return &addNoteAlertResp, nil
}

func (cli *OpsGenieAlertClient) ExecuteAction(req alerts.ExecuteActionAlertRequest) (*alerts.ExecuteActionAlertResponse, error) {
	req.ApiKey = cli.apiKey
	// validate the mandatory parameters: apiKey, alertId/alias, action
	if req.ApiKey == "" {
		return nil, errors.New("ApiKey is a mandatory field and can not be empty.")
	}
	if req.Action == "" {
		return nil, errors.New("Action is a mandatory field and can not be empty.")
	}
	if req.AlertId == "" && req.Alias == "" {
		return nil, errors.New("At least either Alert Id or Alias should be set in the request.")
	}
	if req.AlertId != "" && req.Alias != "" {
		return nil, errors.New("Either Alert Id or Alias should be set in the request not both.")
	}
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("POST", EXECUTE_ACTION_ALERT_URL, req).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Can not execute the action, unable to send the request.")
	}
	// check the returning HTTP status code
	httpStatusCode := resp.StatusCode
	if httpStatusCode >= 400 && httpStatusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d returned", httpStatusCode))
	}
	if httpStatusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d returned", httpStatusCode))
	}
	var executeActionAlertResp alerts.ExecuteActionAlertResponse
	// check if the response can be unmarshalled
	if err = resp.Body.FromJsonTo(&executeActionAlertResp); err != nil {
		return nil, errors.New("Server response can not be parsed.")
	}
	return &executeActionAlertResp, nil
}

func (cli *OpsGenieAlertClient) AttachFile(req alerts.AttachFileAlertRequest) (*alerts.AttachFileAlertResponse, error) {
	req.ApiKey = cli.apiKey
	// validate the mandatory parameters: apiKey, alertId/alias, attachment
	if req.ApiKey == "" {
		return nil, errors.New("ApiKey is a mandatory field and can not be empty.")
	}
	if req.Attachment == "" {
		return nil, errors.New("Attachment is a mandatory field and can not be empty.")
	}
	if req.AlertId == "" && req.Alias == "" {
		return nil, errors.New("At least either Alert Id or Alias should be set in the request.")
	}
	if req.AlertId != "" && req.Alias != "" {
		return nil, errors.New("Either Alert Id or Alias should be set in the request not both.")
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fileName := req.Attachment
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.New("Attachment can not be opened for reading.")
	}
	// add the attachment
	fw, err := w.CreateFormFile("attachment", fileName)
	if err != nil {
		return nil, errors.New("Can not build the request with the field attachment.")
	}
	if _, err := io.Copy(fw, file); err != nil {
		return nil, errors.New("Can not copy the attachment into the request.")
	}

	defer file.Close()

	// Add the other fields
	// empty fields should not be placed into the request
	// otherwise it yields an incomplete boundary exception
	if req.ApiKey != "" {
		if fw, err = w.CreateFormField("apiKey"); err != nil {
			return nil, errors.New("Can not build the request with the field apiKey.")
		}
		if _, err = fw.Write([]byte(req.ApiKey)); err != nil {
			return nil, errors.New("Can not write the ApiKey value into the request.")
		}
	}
	if req.AlertId != "" {
		if fw, err = w.CreateFormField("alertId"); err != nil {
			return nil, errors.New("Can not build the request with the field alertId.")
		}
		if _, err = fw.Write([]byte(req.AlertId)); err != nil {
			return nil, errors.New("Can not write the AlertId value into the request.")
		}
	}
	if req.Alias != "" {
		if fw, err = w.CreateFormField("alias"); err != nil {
			return nil, errors.New("Can not build the request with the field alias.")
		}
		if _, err = fw.Write([]byte(req.Alias)); err != nil {
			return nil, errors.New("Can not write the Alias value into the request.")
		}
	}
	if req.User != "" {
		if fw, err = w.CreateFormField("user"); err != nil {
			return nil, errors.New("Can not build the request with the field user.")
		}
		if _, err = fw.Write([]byte(req.User)); err != nil {
			return nil, errors.New("Can not write the User value into the request.")
		}
	}
	if req.Source != "" {
		if fw, err = w.CreateFormField("source"); err != nil {
			return nil, errors.New("Can not build the request with the field source.")
		}
		if _, err = fw.Write([]byte(req.Source)); err != nil {
			return nil, errors.New("Can not write the Source value into the request.")
		}
	}
	if req.IndexFile != "" {
		if fw, err = w.CreateFormField("indexFile"); err != nil {
			return nil, errors.New("Can not build the request with the field indexFile.")
		}
		if _, err = fw.Write([]byte(req.IndexFile)); err != nil {
			return nil, errors.New("Can not write the IndexFile value into the request.")
		}
	}
	if req.Note != "" {
		if fw, err = w.CreateFormField("note"); err != nil {
			return nil, errors.New("Can not build the request with the field note.")
		}
		if _, err = fw.Write([]byte(req.Note)); err != nil {
			return nil, errors.New("Can not write the Note value into the request.")
		}
	}

	w.Close()

	httpReq, err := http.NewRequest("POST", ATTACH_FILE_ALERT_URL, &b)
	if err != nil {
		return nil, errors.New("Can not create the multipart/form-data request.")
	}
	httpReq.Header.Set("Content-Type", w.FormDataContentType())
	client := &http.Client{}
	// proxy settings
	if cli.proxy != "" {
		proxyUrl, proxyErr := url.Parse(cli.proxy)
		if proxyErr != nil {
			return nil, errors.New("Can not set the proxy configuration")
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}
	var res *http.Response
	for i := 0; i < cli.retries; i++ {
		res, err = client.Do(httpReq)
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}

	if err != nil {
		return nil, errors.New("Can not attach the file, unable to send the request.")
	}

	// check the returning HTTP status code
	httpStatusCode := res.StatusCode
	if httpStatusCode >= 400 && httpStatusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d returned", httpStatusCode))
	}
	if httpStatusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d returned", httpStatusCode))
	}

	attachFileAlertResp := alerts.AttachFileAlertResponse{Status: res.Status, Code: res.StatusCode}
	return &attachFileAlertResp, nil
}
