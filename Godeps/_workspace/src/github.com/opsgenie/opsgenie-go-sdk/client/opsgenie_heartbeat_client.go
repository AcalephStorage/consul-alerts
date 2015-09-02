package client

import (
	"errors"
	"fmt"
	goreq "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/franela/goreq"
	goquery "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/google/go-querystring/query"
	heartbeat "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/opsgenie/opsgenie-go-sdk/heartbeat"
	"time"
)

const (
	ADD_HEARTBEAT_URL     = ENDPOINT_URL + "/v1/json/heartbeat"
	UPDATE_HEARTBEAT_URL  = ENDPOINT_URL + "/v1/json/heartbeat"
	ENABLE_HEARTBEAT_URL  = ENDPOINT_URL + "/v1/json/heartbeat/enable"
	DISABLE_HEARTBEAT_URL = ENDPOINT_URL + "/v1/json/heartbeat/disable"
	DELETE_HEARTBEAT_URL  = ENDPOINT_URL + "/v1/json/heartbeat"
	GET_HEARTBEAT_URL     = ENDPOINT_URL + "/v1/json/heartbeat"
	LIST_HEARTBEAT_URL    = ENDPOINT_URL + "/v1/json/heartbeat"
	SEND_HEARTBEAT_URL    = ENDPOINT_URL + "/v1/json/heartbeat/send"
)

type OpsGenieHeartbeatClient struct {
	apiKey  string
	proxy   string
	retries int
}

func (cli *OpsGenieHeartbeatClient) buildRequest(method string, uri string, body interface{}) goreq.Request {
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

func (cli *OpsGenieHeartbeatClient) SetConnectionTimeout(timeoutInSeconds time.Duration) {
	goreq.SetConnectTimeout(timeoutInSeconds * time.Second)
}

func (cli *OpsGenieHeartbeatClient) SetMaxRetryAttempts(retries int) {
	cli.retries = retries
}

func (cli *OpsGenieHeartbeatClient) Add(req heartbeat.AddHeartbeatRequest) (*heartbeat.AddHeartbeatResponse, error) {
	req.ApiKey = cli.apiKey
	// validate mandatory fields: apiKey, name
	if req.ApiKey == "" {
		return nil, errors.New("Api Key is a mandatory field and can not be empty")
	}
	if req.Name == "" {
		return nil, errors.New("Heart beat name is a mandatory field and can not be empty")
	}
	// send the request
	resp, err := cli.buildRequest("POST", ADD_HEARTBEAT_URL, req).Do()
	// resp, err := goreq.Request{ Method: "POST", Uri: ADD_HEARTBEAT_URL, Body: req, }.Do()
	if err != nil {
		return nil, errors.New("Can not add a new heart beat, unable to send the request")
	}
	// check for the returning http status, 4xx: client errors, 5xx: server errors
	statusCode := resp.StatusCode
	if statusCode >= 400 && statusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d occured", statusCode))
	}
	if statusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d occured", statusCode))
	}
	// try to parse the returning JSON into the response
	var addHeartbeatResp heartbeat.AddHeartbeatResponse
	if err = resp.Body.FromJsonTo(&addHeartbeatResp); err != nil {
		return nil, errors.New("Server response can not be parsed")
	}
	// parsed successfuly with no errors
	return &addHeartbeatResp, nil
}

// Update Heartbeat is used to change configuration of existing heartbeats.
// Mandatory Parameters:
// 	- id: 		Id of the heartbeat
// 	- apiKey: 	API key is used for authenticating API requests
// Optional Parameters
// 	- name: 			Name of the heartbeat
// 	- interval: 		Specifies how often a heartbeat message should be expected.
// 	- intervalUnit: 	interval specified as minutes, hours or days
// 	- description:	 	An optional description of the heartbeat
// 	- enabled: 			Enable/disable heartbeat monitoring
func (cli *OpsGenieHeartbeatClient) Update(req heartbeat.UpdateHeartbeatRequest) (*heartbeat.UpdateHeartbeatResponse, error) {
	req.ApiKey = cli.apiKey
	// validate mandatory fields: apiKey, id
	if req.ApiKey == "" {
		return nil, errors.New("Api Key is a mandatory field and can not be empty")
	}
	if req.Id == "" {
		return nil, errors.New("Id is a mandatory field and can not be empty")
	}
	// send the request
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("POST", UPDATE_HEARTBEAT_URL, req).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Can not update the heartbeat, unable to send the request")
	}
	// check for the returning http status, 4xx: client errors, 5xx: server errors
	statusCode := resp.StatusCode
	if statusCode >= 400 && statusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d occured", statusCode))
	}
	if statusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d occured", statusCode))
	}
	// try to parse the returning JSON into the response
	var updateHeartbeatResp heartbeat.UpdateHeartbeatResponse
	if err = resp.Body.FromJsonTo(&updateHeartbeatResp); err != nil {
		return nil, errors.New("Server response can not be parsed")
	}
	// parsed successfuly with no errors
	return &updateHeartbeatResp, nil
}

func (cli *OpsGenieHeartbeatClient) Enable(req heartbeat.EnableHeartbeatRequest) (*heartbeat.EnableHeartbeatResponse, error) {
	req.ApiKey = cli.apiKey
	// validate mandatory fields: apiKey, name/id
	if req.ApiKey == "" {
		return nil, errors.New("Api Key is a mandatory field and can not be empty")
	}
	if req.Name == "" && req.Id == "" {
		return nil, errors.New("One of the 'Name' and 'Id' parameters should be supplied at least")
	}
	if req.Name != "" && req.Id != "" {
		return nil, errors.New("Either 'Name' or 'Id' field should be supplied not both")
	}
	// send the request in a query string
	v, _ := goquery.Values(req)
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("POST", ENABLE_HEARTBEAT_URL+"?"+v.Encode(), nil).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Can not enable the heart beat, unable to send the request")
	}
	// check for the returning http status, 4xx: client errors, 5xx: server errors
	statusCode := resp.StatusCode
	if statusCode >= 400 && statusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d occured", statusCode))
	}
	if statusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d occured", statusCode))
	}
	// try to parse the returning JSON into the response
	var enableHeartbeatResp heartbeat.EnableHeartbeatResponse
	if err = resp.Body.FromJsonTo(&enableHeartbeatResp); err != nil {
		return nil, errors.New("Server response can not be parsed")
	}
	// parsed successfuly with no errors
	return &enableHeartbeatResp, nil
}

func (cli *OpsGenieHeartbeatClient) Disable(req heartbeat.DisableHeartbeatRequest) (*heartbeat.DisableHeartbeatResponse, error) {
	req.ApiKey = cli.apiKey
	// validate mandatory fields: apiKey, name, id
	if req.ApiKey == "" {
		return nil, errors.New("Api Key is a mandatory field and can not be empty")
	}
	if req.Name == "" && req.Id == "" {
		return nil, errors.New("One of the 'Name' and 'Id' parameters should be supplied at least")
	}
	if req.Name != "" && req.Id != "" {
		return nil, errors.New("Either 'Name' or 'Id' field should be supplied not both")
	}
	// send the request in a query string
	v, _ := goquery.Values(req)
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("POST", DISABLE_HEARTBEAT_URL+"?"+v.Encode(), nil).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Can not disable the heart beat, unable to send the request")
	}
	// check for the returning http status, 4xx: client errors, 5xx: server errors
	statusCode := resp.StatusCode
	if statusCode >= 400 && statusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d occured", statusCode))
	}
	if statusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d occured", statusCode))
	}
	// try to parse the returning JSON into the response
	var disableHeartbeatResp heartbeat.DisableHeartbeatResponse
	if err = resp.Body.FromJsonTo(&disableHeartbeatResp); err != nil {
		return nil, errors.New("Server response can not be parsed")
	}
	// parsed successfuly with no errors
	return &disableHeartbeatResp, nil
}

func (cli *OpsGenieHeartbeatClient) Delete(req heartbeat.DeleteHeartbeatRequest) (*heartbeat.DeleteHeartbeatResponse, error) {
	req.ApiKey = cli.apiKey
	// validate mandatory fields: apiKey, name/id
	if req.ApiKey == "" {
		return nil, errors.New("Api Key is a mandatory field and can not be empty")
	}
	if req.Name == "" && req.Id == "" {
		return nil, errors.New("Either Name or Id field is mandatory and can not be empty")
	}
	if req.Name != "" && req.Id != "" {
		return nil, errors.New("Either Name or Id field should be supplied not both")
	}
	// send the request in a query string
	v, _ := goquery.Values(req)
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("DELETE", DELETE_HEARTBEAT_URL+"?"+v.Encode(), nil).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}

	// resp, err := goreq.Request{ Method: "DELETE", Uri: DELETE_HEARTBEAT_URL + "?" + v.Encode(), }.Do()
	if err != nil {
		return nil, errors.New("Can not delete the heart beat, unable to send the request")
	}
	// check for the returning http status, 4xx: client errors, 5xx: server errors
	statusCode := resp.StatusCode
	if statusCode >= 400 && statusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d occured", statusCode))
	}
	if statusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d occured", statusCode))
	}
	// try to parse the returning JSON into the response
	var deleteHeartbeatResp heartbeat.DeleteHeartbeatResponse
	if err = resp.Body.FromJsonTo(&deleteHeartbeatResp); err != nil {
		return nil, errors.New("Server response can not be parsed")
	}
	// parsed successfuly with no errors
	return &deleteHeartbeatResp, nil
}

func (cli *OpsGenieHeartbeatClient) Get(req heartbeat.GetHeartbeatRequest) (*heartbeat.GetHeartbeatResponse, error) {
	req.ApiKey = cli.apiKey
	// validate mandatory fields: apiKey, name/id
	if req.ApiKey == "" {
		return nil, errors.New("Api Key is a mandatory field and can not be empty")
	}
	if req.Name == "" && req.Id == "" {
		return nil, errors.New("One of the 'Name' and 'Id' parameters should be supplied at least")
	}
	if req.Name != "" && req.Id != "" {
		return nil, errors.New("Either 'Name' or 'Id' field should be supplied not both")
	}
	// send the request in a query string
	v, _ := goquery.Values(req)
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("GET", GET_HEARTBEAT_URL+"?"+v.Encode(), nil).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Can not retrieve the heart beat, unable to send the request")
	}
	// check for the returning http status, 4xx: client errors, 5xx: server errors
	statusCode := resp.StatusCode
	if statusCode >= 400 && statusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d occured", statusCode))
	}
	if statusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d occured", statusCode))
	}
	// try to parse the returning JSON into the response
	var getHeartbeatResp heartbeat.GetHeartbeatResponse
	if err = resp.Body.FromJsonTo(&getHeartbeatResp); err != nil {
		return nil, errors.New("Server response can not be parsed")
	}
	// parsed successfuly with no errors
	return &getHeartbeatResp, nil
}

func (cli *OpsGenieHeartbeatClient) List(req heartbeat.ListHeartbeatsRequest) (*heartbeat.ListHeartbeatsResponse, error) {
	req.ApiKey = cli.apiKey
	// validate the mandatory field: apiKey
	if req.ApiKey == "" {
		return nil, errors.New("Api Key is a mandatory field and can not be empty")
	}
	// send the request in a query string
	v, _ := goquery.Values(req)
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("GET", LIST_HEARTBEAT_URL+"?"+v.Encode(), nil).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Can not retrieve the list of heart beats, unable to send the request")
	}
	// check for the returning http status, 4xx: client errors, 5xx: server errors
	statusCode := resp.StatusCode
	if statusCode >= 400 && statusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d occured", statusCode))
	}
	if statusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d occured", statusCode))
	}
	// try to parse the returning JSON into the response
	var listHeartbeatsResp heartbeat.ListHeartbeatsResponse
	if err = resp.Body.FromJsonTo(&listHeartbeatsResp); err != nil {
		return nil, errors.New("Server response can not be parsed")
	}
	// parsed successfuly with no errors
	return &listHeartbeatsResp, nil
}

func (cli *OpsGenieHeartbeatClient) Send(req heartbeat.SendHeartbeatRequest) (*heartbeat.SendHeartbeatResponse, error) {
	req.ApiKey = cli.apiKey
	// validate the mandatory field: apiKey
	if req.ApiKey == "" {
		return nil, errors.New("Api Key is a mandatory field and can not be empty")
	}
	// send the payload in a POST request
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("POST", SEND_HEARTBEAT_URL, req).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Can not send the heart beat, unable to send the request")
	}
	// check for the returning http status, 4xx: client errors, 5xx: server errors
	statusCode := resp.StatusCode
	if statusCode >= 400 && statusCode < 500 {
		return nil, errors.New(fmt.Sprintf("Client error %d occured", statusCode))
	}
	if statusCode >= 500 {
		return nil, errors.New(fmt.Sprintf("Server error %d occured", statusCode))
	}
	// try to parse the returning JSON into the response
	var sendHeartbeatResp heartbeat.SendHeartbeatResponse
	if err = resp.Body.FromJsonTo(&sendHeartbeatResp); err != nil {
		return nil, errors.New("Server response can not be parsed")
	}
	// parsed successfuly with no errors
	return &sendHeartbeatResp, nil
}
