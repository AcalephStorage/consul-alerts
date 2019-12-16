package client

import (
	"github.com/uchiru/consul-alerts/Godeps/_workspace/src/github.com/opsgenie/opsgenie-go-sdk/schedulev2"
)

// OpsGenieScheduleOverrideV2Client is the data type to make Schedule rule API requests.
type OpsGenieScheduleOverrideV2Client struct {
	RestClient
}

// SetOpsGenieClient sets the embedded OpsGenieClient type of the OpsGenieScheduleOverrideV2Client.
func (cli *OpsGenieScheduleOverrideV2Client) SetOpsGenieClient(ogCli OpsGenieClient) {
	cli.OpsGenieClient = ogCli
}

// Creates method creates a notification rule at OpsGenie.
func (cli *OpsGenieScheduleOverrideV2Client) Create(req schedulev2.CreateScheduleOverrideRequest) (
	*schedulev2.CreateScheduleOverrideResponse,
	error,
) {
	var response schedulev2.CreateScheduleOverrideResponse
	err := cli.sendPostRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Get method returns a schedule from OpsGenie.
func (cli *OpsGenieScheduleOverrideV2Client) Get(req schedulev2.GetScheduleOverrideRequest) (
	*schedulev2.GetScheduleOverrideResponse,
	error,
) {
	var response schedulev2.GetScheduleOverrideResponse
	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Update method updates specified schedule rule.
func (cli *OpsGenieScheduleOverrideV2Client) Update(req schedulev2.UpdateScheduleOverrideRequest) (
	*schedulev2.UpdateScheduleOverrideResponse,
	error,
) {
	var response schedulev2.UpdateScheduleOverrideResponse
	err := cli.sendPutRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Delete method deletes specified schedule rule.
func (cli *OpsGenieScheduleOverrideV2Client) Delete(req schedulev2.DeleteScheduleOverrideRequest) (
	*schedulev2.DeleteScheduleOverrideResponse,
	error,
) {
	var response schedulev2.DeleteScheduleOverrideResponse
	err := cli.sendDeleteRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// List method returns list of existing schedules.
func (cli *OpsGenieScheduleOverrideV2Client) List(req schedulev2.ListScheduleOverrideRequest) (
	*schedulev2.ListScheduleOverrideResponse,
	error,
) {
	var response schedulev2.ListScheduleOverrideResponse
	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
