package alertsv2

import "net/url"

type AddAlertAttachmentRequest struct {
	*AttachmentAlertIdentifier
	AttachmentFilePath    string `json:"alertfile,omitempty"`
	AttachmentFileContent []byte
	AttachmentFileName    string
	User                  string `json:"user,omitempty"`
	IndexFile             string `json:"indexFile,omitempty"`
	ApiKey                string `json:"-"`
}

func (r *AddAlertAttachmentRequest) GenerateUrl() (string, url.Values, error) {
	path, params, err := r.AttachmentAlertIdentifier.GenerateUrl()
	return path + "/attachments", params, err
}

func (r *AddAlertAttachmentRequest) GetApiKey() string {
	return r.ApiKey
}
