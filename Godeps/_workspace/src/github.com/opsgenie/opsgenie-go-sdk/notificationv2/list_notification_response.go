package notificationv2

// ListNotificationResponse is a response of getting notification rules list.
type ListNotificationResponse struct {
	ResponseMeta
	Notifications []Notification `json:"data"`
}
