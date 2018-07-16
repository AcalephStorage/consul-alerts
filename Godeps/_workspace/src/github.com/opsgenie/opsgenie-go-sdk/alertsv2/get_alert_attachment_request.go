package alertsv2

import (
	"net/url"
)

type GetAlertAttachmentRequest struct {
	*AttachmentAlertIdentifier
	AttachmentId string
	ApiKey       string
}

func (r *GetAlertAttachmentRequest) GetApiKey() string {
	return r.ApiKey
}

func (r *GetAlertAttachmentRequest) GenerateUrl() (string, url.Values, error) {
	path, params, err := r.AttachmentAlertIdentifier.GenerateUrl()

	if err != nil {
		return "", nil, err
	}

	return path + "/attachments/" + r.AttachmentId, params, nil
}
