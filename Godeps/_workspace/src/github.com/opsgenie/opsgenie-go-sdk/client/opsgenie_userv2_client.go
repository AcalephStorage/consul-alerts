package client

import (
	"github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/opsgenie/opsgenie-go-sdk/userv2"
)

// OpsGenieUserV2Client is the data type to make User API requests.
type OpsGenieUserV2Client struct {
	RestClient
}

// SetOpsGenieClient sets the embedded OpsGenieClient type of the OpsGenieUserV2Client.
func (cli *OpsGenieUserV2Client) SetOpsGenieClient(ogCli OpsGenieClient) {
	cli.OpsGenieClient = ogCli
}

// List method retrieves the list of users from OpsGenie.
func (cli *OpsGenieUserV2Client) List(req userv2.ListUsersRequest) (*userv2.ListUsersResponse, error) {
	var response userv2.ListUsersResponse
	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// Create method creates new user in OpsGenie.
func (cli *OpsGenieUserV2Client) Create(req userv2.CreateUserRequest) (*userv2.CreateUserResponse, error) {
	var response userv2.CreateUserResponse
	err := cli.sendPostRequest(&req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// Get method returns user.
func (cli *OpsGenieUserV2Client) Get(req userv2.GetUserRequest) (*userv2.GetUserResponse, error) {
	var response userv2.GetUserResponse
	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// Update method updates data of user
func (cli *OpsGenieUserV2Client) Update(req userv2.UpdateUserRequest) (*userv2.UpdateUserResponse, error) {
	var response userv2.UpdateUserResponse
	err := cli.sendPatchRequest(&req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// Delete method deletes user.
func (cli *OpsGenieUserV2Client) Delete(req userv2.DeleteUserRequest) (*userv2.DeleteUserResponse, error) {
	var response userv2.DeleteUserResponse
	err := cli.sendDeleteRequest(&req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// ListEscalations method returns list of user escalations.
func (cli *OpsGenieUserV2Client) ListEscalations(req userv2.ListUserEscalationsRequest) (*userv2.ListUserEscalationsResponse, error) {
	var response userv2.ListUserEscalationsResponse
	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// ListTeams method returns list of user teams.
func (cli *OpsGenieUserV2Client) ListTeams(req userv2.ListUserTeamsRequest) (*userv2.ListUserTeamsResponse, error) {
	var response userv2.ListUserTeamsResponse
	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// ListForwardingRules method returns list of user forwarding rules.
func (cli *OpsGenieUserV2Client) ListForwardingRules(req userv2.ListUserForwardingRulesRequest) (*userv2.ListUserForwardingRulesResponse, error) {
	var response userv2.ListUserForwardingRulesResponse
	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// ListSchedules method returns list of user schedules.
func (cli *OpsGenieUserV2Client) ListSchedules(req userv2.ListUserSchedulesRequest) (*userv2.ListUserSchedulesResponse, error) {
	var response userv2.ListUserSchedulesResponse
	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
