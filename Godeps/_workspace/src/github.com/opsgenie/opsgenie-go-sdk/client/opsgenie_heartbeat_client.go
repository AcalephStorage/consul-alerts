package client

import (
	"errors"
	"github.com/uchiru/consul-alerts/Godeps/_workspace/src/github.com/opsgenie/opsgenie-go-sdk/heartbeat"
	"github.com/uchiru/consul-alerts/Godeps/_workspace/src/github.com/opsgenie/opsgenie-go-sdk/logging"
	"time"
)

const (
	listHeartbeatURL = "/v1/json/heartbeat"
	sendHeartbeatURL = "/v1/json/heartbeat/send"
)

// OpsGenieHeartbeatClient is the data type to make Heartbeat API requests.
type OpsGenieHeartbeatClient struct {
	RestClient
}

// SetOpsGenieClient sets the embedded OpsGenieClient type of the OpsGenieHeartbeatClient.
func (cli *OpsGenieHeartbeatClient) SetOpsGenieClient(ogCli OpsGenieClient) {
	cli.OpsGenieClient = ogCli
}

// Add method creates a heartbeat at OpsGenie.
func (cli *OpsGenieHeartbeatClient) Add(req heartbeat.AddHeartbeatRequest) (*heartbeat.AddHeartbeatResponse, error) {
	var response heartbeat.HeartbeatResponseV2

	err := cli.sendPostRequest(&req, &response)

	if err != nil {
		return nil, err
	}
	result := convertAddResponseToV1Response(&response)
	result.Code = 201
	return result, nil
}

// Update method changes configuration of an existing heartbeat at OpsGenie.
func (cli *OpsGenieHeartbeatClient) Update(req heartbeat.UpdateHeartbeatRequest) (*heartbeat.UpdateHeartbeatResponse, error) {
	var response heartbeat.HeartbeatMetaResponseV2
	err := cli.sendPatchRequest(&req, &response)

	if err != nil {
		return nil, err
	}

	return convertUpdateToV1Response(&response), nil
}

// Enable method enables an heartbeat at OpsGenie.
func (cli *OpsGenieHeartbeatClient) Enable(req heartbeat.EnableHeartbeatRequest) (*heartbeat.EnableHeartbeatResponse, error) {
	var response heartbeat.HeartbeatMetaResponseV2

	err := cli.sendPostRequest(&req, &response)

	if err != nil {
		return nil, err
	}

	var result heartbeat.EnableHeartbeatResponse
	result.Status = "successful"
	result.Code = 200
	return &result, nil
}

// Disable method disables an heartbeat at OpsGenie.
func (cli *OpsGenieHeartbeatClient) Disable(req heartbeat.DisableHeartbeatRequest) (*heartbeat.DisableHeartbeatResponse, error) {
	var response heartbeat.HeartbeatMetaResponseV2

	err := cli.sendPostRequest(&req, &response)

	if err != nil {
		return nil, err
	}

	var result heartbeat.DisableHeartbeatResponse
	result.Status = "successful"
	result.Code = 200
	return &result, nil
}

// Delete method deletes an heartbeat from OpsGenie.
func (cli *OpsGenieHeartbeatClient) Delete(req heartbeat.DeleteHeartbeatRequest) (*heartbeat.DeleteHeartbeatResponse, error) {
	var response heartbeat.HeartbeatMetaResponseV2

	err := cli.sendDeleteRequest(&req, &response)

	if err != nil {
		return nil, err
	}

	var result heartbeat.DeleteHeartbeatResponse
	result.Status = "Deleted"
	result.Code = 200
	return &result, nil
}

// Get method retrieves an heartbeat with details from OpsGenie.
func (cli *OpsGenieHeartbeatClient) Get(req heartbeat.GetHeartbeatRequest) (*heartbeat.GetHeartbeatResponse, error) {
	var response heartbeat.HeartbeatResponseV2

	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return convertGetToV1Response(&response), nil
}

// Deprecated: List method retrieves heartbeats from OpsGenie.
func (cli *OpsGenieHeartbeatClient) List(req heartbeat.ListHeartbeatsRequest) (*heartbeat.ListHeartbeatsResponse, error) {
	req.APIKey = cli.apiKey
	resp, err := cli.sendRequest(cli.buildGetRequest(listHeartbeatURL, req))

	if resp == nil {
		return nil, err
	}
	defer resp.Body.Close()

	var listHeartbeatsResp heartbeat.ListHeartbeatsResponse
	if err = resp.Body.FromJsonTo(&listHeartbeatsResp); err != nil {
		message := "Server response can not be parsed, " + err.Error()
		logging.Logger().Warn(message)
		return nil, errors.New(message)
	}

	return &listHeartbeatsResp, nil
}

// Deprecated: Send method sends an Heartbeat Signal to OpsGenie.
func (cli *OpsGenieHeartbeatClient) Send(req heartbeat.SendHeartbeatRequest) (*heartbeat.SendHeartbeatResponse, error) {
	req.APIKey = cli.apiKey
	resp, err := cli.sendRequest(cli.buildPostRequest(sendHeartbeatURL, req))

	if resp == nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sendHeartbeatResp heartbeat.SendHeartbeatResponse
	if err = resp.Body.FromJsonTo(&sendHeartbeatResp); err != nil {
		message := "Server response can not be parsed, " + err.Error()
		logging.Logger().Warn(message)
		return nil, errors.New(message)
	}

	return &sendHeartbeatResp, nil
}

// Send method sends an Heartbeat Signal to OpsGenie.
func (cli *OpsGenieHeartbeatClient) Ping(req heartbeat.PingHeartbeatRequest) (*AsyncRequestResponse, error) {
	return cli.sendAsyncPostRequest(&req)
}

func convertAddResponseToV1Response(responseV2 *heartbeat.HeartbeatResponseV2) *heartbeat.AddHeartbeatResponse {
	data := responseV2.Data

	var result = heartbeat.AddHeartbeatResponse{}

	result.Name = data.Name
	result.Status = "successful"

	return &result
}

func convertUpdateToV1Response(responseV2 *heartbeat.HeartbeatMetaResponseV2) *heartbeat.UpdateHeartbeatResponse {
	data := responseV2.Data

	var result = heartbeat.UpdateHeartbeatResponse{}

	result.Name = data.Name
	result.Status = "successful"
	result.Code = 200
	return &result
}

func convertGetToV1Response(responseV2 *heartbeat.HeartbeatResponseV2) *heartbeat.GetHeartbeatResponse {
	data := responseV2.Data

	var result = heartbeat.GetHeartbeatResponse{}

	result.Name = data.Name
	result.Description = data.Description
	result.Interval = data.Interval
	result.IntervalUnit = data.IntervalUnit
	result.Enabled = data.Enabled
	result.Expired = data.Expired

	if data.Expired {
		result.Status = "Expired"
	} else {
		result.Status = "Active"
	}

	if !data.LastHeartbeat.IsZero() {
		result.LastHeartbeat = uint64(data.LastHeartbeat.UnixNano() / int64(time.Millisecond))
	}

	return &result
}
