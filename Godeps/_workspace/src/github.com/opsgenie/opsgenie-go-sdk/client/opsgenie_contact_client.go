package client

import (
	"errors"

	"github.com/uchiru/consul-alerts/Godeps/_workspace/src/github.com/opsgenie/opsgenie-go-sdk/contact"
	"github.com/uchiru/consul-alerts/Godeps/_workspace/src/github.com/opsgenie/opsgenie-go-sdk/logging"
)

const (
	contactURL = "/v1/json/user/contact"
)

// OpsGenieContactClient is the data type to make Contact API requests.
type OpsGenieContactClient struct {
	OpsGenieClient
}

// SetOpsGenieClient sets the embedded OpsGenieClient type of the OpsGenieContactClient.
func (cli *OpsGenieContactClient) SetOpsGenieClient(ogCli OpsGenieClient) {
	cli.OpsGenieClient = ogCli
}

// Create method creates a contact at OpsGenie.
func (cli *OpsGenieContactClient) Create(req contact.CreateContactRequest) (*contact.CreateContactResponse, error) {
	req.APIKey = cli.apiKey
	resp, err := cli.sendRequest(cli.buildPostRequest(contactURL, req))

	if resp == nil {
		return nil, err
	}
	defer resp.Body.Close()

	var createContactResp contact.CreateContactResponse

	if err = resp.Body.FromJsonTo(&createContactResp); err != nil {
		message := "Server response can not be parsed, " + err.Error()
		logging.Logger().Warn(message)
		return nil, errors.New(message)
	}
	return &createContactResp, nil
}

// Delete method deletes a contact at OpsGenie.
func (cli *OpsGenieContactClient) Delete(req contact.DeleteContactRequest) (*contact.DeleteContactResponse, error) {
	req.APIKey = cli.apiKey
	resp, err := cli.sendRequest(cli.buildDeleteRequest(contactURL, req))

	if resp == nil {
		return nil, err
	}
	defer resp.Body.Close()

	var deleteContactResp contact.DeleteContactResponse

	if err = resp.Body.FromJsonTo(&deleteContactResp); err != nil {
		message := "Server response can not be parsed, " + err.Error()
		logging.Logger().Warn(message)
		return nil, errors.New(message)
	}
	return &deleteContactResp, nil
}

// Disable method disables a contact at OpsGenie.
func (cli *OpsGenieContactClient) Disable(req contact.DisableContactRequest) (*contact.DisableContactResponse, error) {
	req.APIKey = cli.apiKey
	resp, err := cli.sendRequest(cli.buildPostRequest(contactURL+"/disable", req))

	if resp == nil {
		return nil, err
	}
	defer resp.Body.Close()

	var disableContactResp contact.DisableContactResponse

	if err = resp.Body.FromJsonTo(&disableContactResp); err != nil {
		message := "Server response can not be parsed, " + err.Error()
		logging.Logger().Warn(message)
		return nil, errors.New(message)
	}
	return &disableContactResp, nil
}

// Enable method enables a contact at OpsGenie.
func (cli *OpsGenieContactClient) Enable(req contact.EnableContactRequest) (*contact.EnableContactResponse, error) {
	req.APIKey = cli.apiKey
	resp, err := cli.sendRequest(cli.buildPostRequest(contactURL+"/enable", req))

	if resp == nil {
		return nil, err
	}
	defer resp.Body.Close()

	var enableContactResp contact.EnableContactResponse

	if err = resp.Body.FromJsonTo(&enableContactResp); err != nil {
		message := "Server response can not be parsed, " + err.Error()
		logging.Logger().Warn(message)
		return nil, errors.New(message)
	}
	return &enableContactResp, nil
}

// Get method retrieves specified contact details from OpsGenie.
func (cli *OpsGenieContactClient) Get(req contact.GetContactRequest) (*contact.GetContactResponse, error) {
	req.APIKey = cli.apiKey
	resp, err := cli.sendRequest(cli.buildGetRequest(contactURL, req))

	if resp == nil {
		return nil, err
	}
	defer resp.Body.Close()

	var getContactResp contact.GetContactResponse

	if err = resp.Body.FromJsonTo(&getContactResp); err != nil {
		message := "Server response can not be parsed, " + err.Error()
		logging.Logger().Warn(message)
		return nil, errors.New(message)
	}
	return &getContactResp, nil
}

// Update method updates a contact at OpsGenie.
func (cli *OpsGenieContactClient) Update(req contact.UpdateContactRequest) (*contact.UpdateContactResponse, error) {
	req.APIKey = cli.apiKey
	resp, err := cli.sendRequest(cli.buildPostRequest(contactURL, req))

	if resp == nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updateContactResp contact.UpdateContactResponse

	if err = resp.Body.FromJsonTo(&updateContactResp); err != nil {
		message := "Server response can not be parsed, " + err.Error()
		logging.Logger().Warn(message)
		return nil, errors.New(message)
	}
	return &updateContactResp, nil
}
