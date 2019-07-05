package alertsv2

import (
	"net/url"
)

type DeleteAlertAttachmentRequest struct {
	*AttachmentAlertIdentifier
	AttachmentId string
	ApiKey       string
}

func (r *DeleteAlertAttachmentRequest) GetApiKey() string {
	return r.ApiKey
}

func (r *DeleteAlertAttachmentRequest) GenerateUrl() (string, url.Values, error) {
	path, params, err := r.AttachmentAlertIdentifier.GenerateUrl()

	if err != nil {
		return "", nil, err
	}

	return path + "/attachments/" + r.AttachmentId, params, nil
}
