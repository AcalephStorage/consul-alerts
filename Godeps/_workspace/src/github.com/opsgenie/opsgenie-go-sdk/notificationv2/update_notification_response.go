package notificationv2

// UpdateNotificationResponse is a response of updating notification rule.
type UpdateNotificationResponse struct {
	ResponseMeta
	Notification Notification `json:"data"`
}
