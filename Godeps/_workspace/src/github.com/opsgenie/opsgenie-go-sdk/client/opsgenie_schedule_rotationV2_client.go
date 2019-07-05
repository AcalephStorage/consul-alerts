package client

import (
	"github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/opsgenie/opsgenie-go-sdk/schedulev2"
)

// OpsGenieScheduleRotationV2Client is the data type to make Schedule rule API requests.
type OpsGenieScheduleRotationV2Client struct {
	RestClient
}

// SetOpsGenieClient sets the embedded OpsGenieClient type of the OpsGenieScheduleRotationV2Client.
func (cli *OpsGenieScheduleRotationV2Client) SetOpsGenieClient(ogCli OpsGenieClient) {
	cli.OpsGenieClient = ogCli
}

// Creates method creates a notification rule at OpsGenie.
func (cli *OpsGenieScheduleRotationV2Client) Create(req schedulev2.CreateScheduleRotationRequest) (
	*schedulev2.CreateScheduleRotationResponse,
	error,
) {
	var response schedulev2.CreateScheduleRotationResponse
	err := cli.sendPostRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Get method returns a schedule from OpsGenie.
func (cli *OpsGenieScheduleRotationV2Client) Get(req schedulev2.GetScheduleRotationRequest) (
	*schedulev2.GetScheduleRotationResponse,
	error,
) {
	var response schedulev2.GetScheduleRotationResponse
	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Update method updates specified schedule rule.
func (cli *OpsGenieScheduleRotationV2Client) Update(req schedulev2.UpdateScheduleRotationRequest) (
	*schedulev2.UpdateScheduleRotationResponse,
	error,
) {
	var response schedulev2.UpdateScheduleRotationResponse
	err := cli.sendPatchRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Delete method deletes specified schedule rule.
func (cli *OpsGenieScheduleRotationV2Client) Delete(req schedulev2.DeleteScheduleRotationRequest) (
	*schedulev2.DeleteScheduleRotationResponse,
	error,
) {
	var response schedulev2.DeleteScheduleRotationResponse
	err := cli.sendDeleteRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// List method returns list of existing schedules.
func (cli *OpsGenieScheduleRotationV2Client) List(req schedulev2.ListScheduleRotationRequest) (
	*schedulev2.ListScheduleRotationResponse,
	error,
) {
	var response schedulev2.ListScheduleRotationResponse
	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
