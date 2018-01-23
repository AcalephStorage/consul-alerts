package notificationv2

const (
	// List of time periods that notification for schedule start/end will be sent.
	JustBeforeNotificationTime        	NotificationTime = "just-before"
	FifteenMinutesAgoNotificationTime 	NotificationTime = "15-minutes-ago"
	OneHourAgoNotificationTime        	NotificationTime = "1-hour-ago"
	OneDayAgoNotificationTime         	NotificationTime = "1-day-ago"
)

// NotificationTime is type of time periods that notification for start/end will be sent.
type NotificationTime string
