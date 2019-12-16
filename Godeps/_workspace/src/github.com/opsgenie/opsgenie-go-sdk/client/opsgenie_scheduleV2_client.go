package client

import (
	"github.com/uchiru/consul-alerts/Godeps/_workspace/src/github.com/opsgenie/opsgenie-go-sdk/schedulev2"
)

// OpsGenieScheduleV2Client is the data type to make Schedule rule API requests.
type OpsGenieScheduleV2Client struct {
	RestClient
}

// SetOpsGenieClient sets the embedded OpsGenieClient type of the OpsGenieScheduleV2Client.
func (cli *OpsGenieScheduleV2Client) SetOpsGenieClient(ogCli OpsGenieClient) {
	cli.OpsGenieClient = ogCli
}

// Creates method creates a notification rule at OpsGenie.
func (cli *OpsGenieScheduleV2Client) Create(req schedulev2.CreateScheduleRequest) (
	*schedulev2.CreateScheduleResponse,
	error,
) {
	var response schedulev2.CreateScheduleResponse
	err := cli.sendPostRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Get method returns a schedule from OpsGenie.
func (cli *OpsGenieScheduleV2Client) Get(req schedulev2.GetScheduleRequest) (
	*schedulev2.GetScheduleResponse,
	error,
) {
	var response schedulev2.GetScheduleResponse
	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Update method updates specified schedule rule.
func (cli *OpsGenieScheduleV2Client) Update(req schedulev2.UpdateScheduleRequest) (
	*schedulev2.UpdateScheduleResponse,
	error,
) {
	var response schedulev2.UpdateScheduleResponse
	err := cli.sendPatchRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Delete method deletes specified schedule rule.
func (cli *OpsGenieScheduleV2Client) Delete(req schedulev2.DeleteScheduleRequest) (
	*schedulev2.DeleteScheduleResponse,
	error,
) {
	var response schedulev2.DeleteScheduleResponse
	err := cli.sendDeleteRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// List method returns list of existing schedules.
func (cli *OpsGenieScheduleV2Client) List(req schedulev2.ListScheduleRequest) (
	*schedulev2.ListScheduleResponse,
	error,
) {
	var response schedulev2.ListScheduleResponse
	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
