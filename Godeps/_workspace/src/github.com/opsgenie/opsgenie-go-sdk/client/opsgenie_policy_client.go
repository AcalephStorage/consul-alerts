package client

import (
	"errors"
	"fmt"
	goreq "github.com/vincentvu/consul-alerts/Godeps/_workspace/src/github.com/franela/goreq"
	policy "github.com/vincentvu/consul-alerts/Godeps/_workspace/src/github.com/opsgenie/opsgenie-go-sdk/policy"
	"time"
)

const (
	ENABLE_POLICY_URL  = ENDPOINT_URL + "/v1/json/policy/enable"
	DISABLE_POLICY_URL = ENDPOINT_URL + "/v1/json/policy/disable"
)

type OpsGeniePolicyClient struct {
	apiKey  string
	proxy   string
	retries int
}

func (cli *OpsGeniePolicyClient) buildRequest(method string, uri string, body interface{}) goreq.Request {
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

func (cli *OpsGeniePolicyClient) SetConnectionTimeout(timeoutInSeconds time.Duration) {
	goreq.SetConnectTimeout(timeoutInSeconds * time.Second)
}

func (cli *OpsGeniePolicyClient) SetMaxRetryAttempts(retries int) {
	cli.retries = retries
}

func (cli *OpsGeniePolicyClient) Enable(req policy.EnablePolicyRequest) (*policy.EnablePolicyResponse, error) {
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
		resp, err = cli.buildRequest("POST", ENABLE_POLICY_URL, req).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Can not enable the policy, unable to send the request")
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
	var enablePolicyResp policy.EnablePolicyResponse
	if err = resp.Body.FromJsonTo(&enablePolicyResp); err != nil {
		return nil, errors.New("Server response can not be parsed")
	}
	// parsed successfuly with no errors
	return &enablePolicyResp, nil
}

func (cli *OpsGeniePolicyClient) Disable(req policy.DisablePolicyRequest) (*policy.DisablePolicyResponse, error) {
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
		resp, err = cli.buildRequest("POST", DISABLE_POLICY_URL, req).Do()
		if err == nil {
			break
		}
		time.Sleep(TIME_SLEEP_BETWEEN_REQUESTS)
	}
	if err != nil {
		return nil, errors.New("Can not disable the policy, unable to send the request")
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
	var disablePolicyResp policy.DisablePolicyResponse
	if err = resp.Body.FromJsonTo(&disablePolicyResp); err != nil {
		return nil, errors.New("Server response can not be parsed")
	}
	// parsed successfuly with no errors
	return &disablePolicyResp, nil
}
