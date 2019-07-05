package alertsv2

import (
	"net/url"
)

type ListAlertAttachmentRequest struct {
	*AttachmentAlertIdentifier
	ApiKey string
}

func (r *ListAlertAttachmentRequest) GetApiKey() string {
	return r.ApiKey
}

func (r *ListAlertAttachmentRequest) GenerateUrl() (string, url.Values, error) {
	path, params, err := r.AttachmentAlertIdentifier.GenerateUrl()

	if err != nil {
		return "", nil, err
	}

	return path + "/attachments", params, nil
}
