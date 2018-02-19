package client

import (
	"errors"
	"fmt"
	goreq "github.com/Difrex/consul-alerts/Godeps/_workspace/src/github.com/franela/goreq"
	integration "github.com/Difrex/consul-alerts/Godeps/_workspace/src/github.com/opsgenie/opsgenie-go-sdk/integration"
	"time"
)

const (
	ENABLE_INTEGRATION_URL  = ENDPOINT_URL + "/v1/json/integration/enable"
	DISABLE_INTEGRATION_URL = ENDPOINT_URL + "/v1/json/integration/disable"
)

type OpsGenieIntegrationClient struct {
	apiKey  string
	proxy   string
	retries int
}

func (cli *OpsGenieIntegrationClient) buildRequest(method string, uri string, body interface{}) goreq.Request {
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

func (cli *OpsGenieIntegrationClient) SetConnectionTimeout(timeoutInSeconds time.Duration) {
	goreq.SetConnectTimeout(timeoutInSeconds * time.Second)
}

func (cli *OpsGenieIntegrationClient) SetMaxRetryAttempts(retries int) {
	cli.retries = retries
}

func (cli *OpsGenieIntegrationClient) Enable(req integration.EnableIntegrationRequest) (*integration.EnableIntegrationResponse, error) {
	req.ApiKey = cli.apiKey
	// validate mandatory fields: id/name, apiKey
	if req.ApiKey == "" && req.Id == "" {
		return nil, errors.New("Api Key or Id should be provided")
	}
	if req.ApiKey != "" && req.Id != "" {
		return nil, errors.New("Either Api Key or Id should be provided, not both")
	}
	// send the request
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("POST", ENABLE_INTEGRATION_URL, req).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Can not enable the integration, unable to send the request")
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
	var enableIntegrationResp integration.EnableIntegrationResponse
	if err = resp.Body.FromJsonTo(&enableIntegrationResp); err != nil {
		return nil, errors.New("Server response can not be parsed")
	}
	// parsed successfuly with no errors
	return &enableIntegrationResp, nil
}

func (cli *OpsGenieIntegrationClient) Disable(req integration.DisableIntegrationRequest) (*integration.DisableIntegrationResponse, error) {
	req.ApiKey = cli.apiKey
	// validate mandatory fields: id/name, apiKey
	if req.ApiKey == "" && req.Id == "" {
		return nil, errors.New("Api Key or Id should be provided")
	}
	if req.ApiKey != "" && req.Id != "" {
		return nil, errors.New("Either Api Key or Id should be provided, not both")
	}
	// send the request
	var resp *goreq.Response
	var err error
	for i := 0; i < cli.retries; i++ {
		resp, err = cli.buildRequest("POST", DISABLE_INTEGRATION_URL, req).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Can not disable the integration, unable to send the request")
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
	var disableIntegrationResp integration.DisableIntegrationResponse
	if err = resp.Body.FromJsonTo(&disableIntegrationResp); err != nil {
		return nil, errors.New("Server response can not be parsed")
	}
	// parsed successfuly with no errors
	return &disableIntegrationResp, nil
}
