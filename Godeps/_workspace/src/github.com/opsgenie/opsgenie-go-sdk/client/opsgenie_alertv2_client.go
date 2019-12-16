package client

import (
	"errors"
	"github.com/uchiru/consul-alerts/Godeps/_workspace/src/github.com/franela/goreq"
	"github.com/uchiru/consul-alerts/Godeps/_workspace/src/github.com/opsgenie/opsgenie-go-sdk/alertsv2"
	"github.com/uchiru/consul-alerts/Godeps/_workspace/src/github.com/opsgenie/opsgenie-go-sdk/alertsv2/savedsearches"
)

// OpsGenieAlertClient is the data type to make Alert API requests.

type OpsGenieAlertV2Client struct {
	RestClient
}

// SetOpsGenieClient sets the embedded OpsGenieClient type of the OpsGenieAlertV2Client.
func (cli *OpsGenieAlertV2Client) SetOpsGenieClient(ogCli OpsGenieClient) {
	cli.OpsGenieClient = ogCli
}

// Returns AlertDetails
func (cli *OpsGenieAlertV2Client) Get(req alertsv2.GetAlertRequest) (*alertsv2.DetailedAlertResponse, error) {
	var response alertsv2.DetailedAlertResponse
	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// Retrieves the alerts from OpsGenie
func (cli *OpsGenieAlertV2Client) List(req alertsv2.ListAlertRequest) (*alertsv2.ListAlertResponse, error) {
	var response alertsv2.ListAlertResponse
	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// Retrieves the alert recipients from OpsGenie
func (cli *OpsGenieAlertV2Client) ListAlertRecipients(req alertsv2.ListAlertRecipientsRequest) (*alertsv2.ListAlertRecipientsResponse, error) {
	var response alertsv2.ListAlertRecipientsResponse
	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// Retrieves the alert logs from OpsGenie
func (cli *OpsGenieAlertV2Client) ListAlertLogs(req alertsv2.ListAlertLogsRequest) (*alertsv2.ListAlertLogsResponse, error) {
	var response alertsv2.ListAlertLogsResponse
	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// Retrieves the alert notes from OpsGenie
func (cli *OpsGenieAlertV2Client) ListAlertNotes(req alertsv2.ListAlertNotesRequest) (*alertsv2.ListAlertNotesResponse, error) {
	var response alertsv2.ListAlertNotesResponse
	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// Sends acknowledge request an alert
func (cli *OpsGenieAlertV2Client) Acknowledge(req alertsv2.AcknowledgeRequest) (*AsyncRequestResponse, error) {
	return cli.sendAsyncPostRequest(&req)
}

// Sends close request for an alert
func (cli *OpsGenieAlertV2Client) Close(req alertsv2.CloseRequest) (*AsyncRequestResponse, error) {
	return cli.sendAsyncPostRequest(&req)
}

// Sends unacknowledge request for an alert
func (cli *OpsGenieAlertV2Client) Unacknowledge(req alertsv2.UnacknowledgeRequest) (*AsyncRequestResponse, error) {
	return cli.sendAsyncPostRequest(&req)
}

// Sends snooze request for an alert
func (cli *OpsGenieAlertV2Client) Snooze(req alertsv2.SnoozeRequest) (*AsyncRequestResponse, error) {
	return cli.sendAsyncPostRequest(&req)
}

// Sends custom execution action for an alert
func (cli *OpsGenieAlertV2Client) ExecuteCustomAction(req alertsv2.ExecuteCustomActionRequest) (*AsyncRequestResponse, error) {
	return cli.sendAsyncPostRequest(&req)
}

// Creates an alert
func (cli *OpsGenieAlertV2Client) Create(req alertsv2.CreateAlertRequest) (*AsyncRequestResponse, error) {
	req.Init()
	return cli.sendAsyncPostRequest(&req)
}

// Deletes an alert
func (cli *OpsGenieAlertV2Client) Delete(req alertsv2.DeleteAlertRequest) (*AsyncRequestResponse, error) {
	var response AsyncRequestResponse

	err := cli.sendDeleteRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Adds a note to the alert
func (cli *OpsGenieAlertV2Client) AddNote(req alertsv2.AddNoteRequest) (*AsyncRequestResponse, error) {
	return cli.sendAsyncPostRequest(&req)
}

// Sends escalation to next request
func (cli *OpsGenieAlertV2Client) EscalateToNext(req alertsv2.EscalateToNextRequest) (*AsyncRequestResponse, error) {
	return cli.sendAsyncPostRequest(&req)
}

// Assigns an alert to a user
func (cli *OpsGenieAlertV2Client) Assign(req alertsv2.AssignAlertRequest) (*AsyncRequestResponse, error) {
	return cli.sendAsyncPostRequest(&req)
}

// Adds a team to the alert
func (cli *OpsGenieAlertV2Client) AddTeamToAlert(req alertsv2.AddTeamToAlertRequest) (*AsyncRequestResponse, error) {
	return cli.sendAsyncPostRequest(&req)
}

// Adds a tag to the alert
func (cli *OpsGenieAlertV2Client) AddTags(req alertsv2.AddTagsToAlertRequest) (*AsyncRequestResponse, error) {
	return cli.sendAsyncPostRequest(&req)
}

// Removes a tag from the alert
func (cli *OpsGenieAlertV2Client) RemoveTags(req alertsv2.RemoveTagsRequest) (*AsyncRequestResponse, error) {
	var response AsyncRequestResponse
	err := cli.sendDeleteRequest(&req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// Adds details to the alert
func (cli *OpsGenieAlertV2Client) AddDetails(req alertsv2.AddDetailsRequest) (*AsyncRequestResponse, error) {
	return cli.sendAsyncPostRequest(&req)
}

// Removes the details from the alert
func (cli *OpsGenieAlertV2Client) RemoveDetails(req alertsv2.RemoveDetailsRequest) (*AsyncRequestResponse, error) {
	var response AsyncRequestResponse
	err := cli.sendDeleteRequest(&req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// Creates a SavedSearch
func (cli *OpsGenieAlertV2Client) CreateSavedSearch(req savedsearches.CreateSavedSearchRequest) (*savedsearches.CreateSavedSearchResponse, error) {
	var response savedsearches.CreateSavedSearchResponse

	err := cli.sendPostRequest(&req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// Retrieves the SavedSearch
func (cli *OpsGenieAlertV2Client) GetSavedSearch(req savedsearches.GetSavedSearchRequest) (*savedsearches.GetSavedSearchResponse, error) {
	var response savedsearches.GetSavedSearchResponse

	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// Updates the SavedSearch
func (cli *OpsGenieAlertV2Client) UpdateSavedSearch(req savedsearches.UpdateSavedSearchRequest) (*savedsearches.UpdateSavedSearchResponse, error) {
	var response savedsearches.UpdateSavedSearchResponse
	err := cli.sendPatchRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Deletes the SavedSearch
func (cli *OpsGenieAlertV2Client) DeleteSavedSearch(req savedsearches.DeleteSavedSearchRequest) (*savedsearches.DeleteSavedSearchResponse, error) {
	var response savedsearches.DeleteSavedSearchResponse
	err := cli.sendDeleteRequest(&req, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// Retrieves list of saved searches
func (cli *OpsGenieAlertV2Client) ListSavedSearches(req alertsv2.LisSavedSearchRequest) (*savedsearches.ListSavedSearchResponse, error) {
	var response savedsearches.ListSavedSearchResponse
	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// Retrieves the status of the async request
func (cli *OpsGenieAlertV2Client) GetAsyncRequestStatus(req alertsv2.GetAsyncRequestStatusRequest) (*alertsv2.GetAsyncRequestStatusResponse, error) {
	var response alertsv2.GetAsyncRequestStatusResponse
	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// adds the attachment file to specified alert
func (cli *OpsGenieAlertV2Client) AttachFile(req alertsv2.AddAlertAttachmentRequest) (*alertsv2.AddAlertAttachmentResponse, error) {
	var response alertsv2.AddAlertAttachmentResponse

	if req.AttachmentFilePath == "" && req.AttachmentFileContent == nil {
		return nil, errors.New("File path or content must be provided.")
	}

	if req.AttachmentFilePath == "" && req.AttachmentFileContent != nil && req.AttachmentFileName == "" {
		return nil, errors.New("File name must be provided if only file content is given.")
	}

	err := cli.sendCreateAttachmentRequest(req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// retrieves the specified attachment with a download link
func (cli *OpsGenieAlertV2Client) GetAttachmentFile(req alertsv2.GetAlertAttachmentRequest) (*alertsv2.GetAlertAttachmentResponse, error) {
	var response alertsv2.GetAlertAttachmentResponse

	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// returns a list of attachment meta information for the specified alert
func (cli *OpsGenieAlertV2Client) ListAlertAttachments(req alertsv2.ListAlertAttachmentRequest) (*alertsv2.ListAlertAttachmentsResponse, error) {
	var response alertsv2.ListAlertAttachmentsResponse

	err := cli.sendGetRequest(&req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// deletes the specified alert attachment
func (cli *OpsGenieAlertV2Client) DeleteAttachment(req alertsv2.DeleteAlertAttachmentRequest) (*alertsv2.DeleteAlertAttachmentResponse, error) {
	var response alertsv2.DeleteAlertAttachmentResponse

	err := cli.sendDeleteRequest(&req, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (cli *RestClient) sendCreateAttachmentRequest(req alertsv2.AddAlertAttachmentRequest, response Response) error {
	path, params, err := req.GenerateUrl()

	if err != nil {
		return err
	}

	path = cli.generateFullPathWithParams(path, params)

	var httpRequest *goreq.Request

	if req.AttachmentFilePath == "" {
		httpRequest, err = cli.buildCreateAttachmentRequestWithBytes(path, req)
	} else {
		httpRequest, err = cli.buildCreateAttachmentRequest(path, req)
	}

	if err != nil {
		return err
	}

	cli.setApiKey(httpRequest, req.GetApiKey())
	httpResponse, err := cli.sendRequest(*httpRequest)

	if err != nil {
		return err
	}

	defer httpResponse.Body.Close()

	err = cli.writeBody(httpResponse, &response)
	if err != nil {
		return err
	}

	cli.setResponseMeta(httpResponse, response)

	return nil
}
