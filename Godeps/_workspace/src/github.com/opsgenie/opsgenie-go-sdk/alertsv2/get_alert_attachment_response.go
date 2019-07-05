package alertsv2

type AttachmentResponse struct {
	Name         string `json:"name,omitempty"`
	DownloadLink string `json:"url,omitempty"`
}

type GetAlertAttachmentResponse struct {
	ResponseMeta
	Attachment AttachmentResponse `json:"data"`
}
