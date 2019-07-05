package alertsv2

type AlertAttachmentMeta struct {
	Name string `json:"name,omitempty"`
	Id   int64  `json:"id,omitempty"`
}

type ListAlertAttachmentsResponse struct {
	ResponseMeta
	AlertAttachments []AlertAttachmentMeta `json:"data"`
}
