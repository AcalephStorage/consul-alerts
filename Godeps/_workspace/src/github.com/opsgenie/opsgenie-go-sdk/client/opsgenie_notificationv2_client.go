package client

import (
	"github.com/uchiru/consul-alerts/Godeps/_workspace/src/github.com/opsgenie/opsgenie-go-sdk/notificationv2"
)

// OpsGenieNotificationV2Client is the data type to make Notification rule API requests.
type OpsGenieNotificationV2Client struct {
	RestClient
}

// SetOpsGenieClient sets the embedded OpsGenieClient type of the OpsGenieNotificationV2Client.
func (cli *OpsGenieNotificationV2Client) SetOpsGenieClient(ogCli OpsGenieClient) {
	cli.OpsGenieClient = ogCli
}

// Creates method creates a notification rule at OpsGenie.
func (cli *OpsGenieNotificationV2Client) Create(req notificationv2.CreateNotificationRequest) (
	*notificationv2.CreateNotificationResponse,
	error,
) {
	var response notificationv2.CreateNotificationResponse
	err := cli.sendPostRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Get method returns a notification from OpsGenie.
func (cli *OpsGenieNotificationV2Client) Get(req notificationv2.GetNotificationRequest) (
	*notificationv2.GetNotificationResponse,
	error,
) {
	var response notificationv2.GetNotificationResponse
	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Update method updates specified notification rule.
func (cli *OpsGenieNotificationV2Client) Update(req notificationv2.UpdateNotificationRequest) (
	*notificationv2.UpdateNotificationResponse,
	error,
) {
	var response notificationv2.UpdateNotificationResponse
	err := cli.sendPatchRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Delete method deletes specified notification rule.
func (cli *OpsGenieNotificationV2Client) Delete(req notificationv2.DeleteNotificationRequest) (
	*notificationv2.DeleteNotificationResponse,
	error,
) {
	var response notificationv2.DeleteNotificationResponse
	err := cli.sendDeleteRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// List method returns list of existing notification rules depending on specified criteria.
func (cli *OpsGenieNotificationV2Client) List(req notificationv2.ListNotificationRequest) (
	*notificationv2.ListNotificationResponse,
	error,
) {
	var response notificationv2.ListNotificationResponse
	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Enable method enables specified notification rule.
func (cli *OpsGenieNotificationV2Client) Enable(req notificationv2.EnableNotificationRequest) (
	*notificationv2.EnableNotificationResponse,
	error,
) {
	var response notificationv2.EnableNotificationResponse
	err := cli.sendPostRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Enable method disables specified notification rule.
func (cli *OpsGenieNotificationV2Client) Disable(req notificationv2.DisableNotificationRequest) (
	*notificationv2.DisableNotificationResponse,
	error,
) {
	var response notificationv2.DisableNotificationResponse
	err := cli.sendPostRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
