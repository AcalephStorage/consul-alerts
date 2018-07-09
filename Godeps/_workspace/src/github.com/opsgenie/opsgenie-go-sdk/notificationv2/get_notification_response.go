package notificationv2

// GetNotificationResponse is a response of getting notification rule.
type GetNotificationResponse struct {
	ResponseMeta
	Notification *Notification `json:"data"`
}
