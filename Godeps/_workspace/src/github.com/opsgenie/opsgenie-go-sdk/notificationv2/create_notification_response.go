package notificationv2

// CreateNotificationResponse is a response of creating alert action.
type CreateNotificationResponse struct {
	ResponseMeta
	Notification Notification `json:"data"`
}
