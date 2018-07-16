package notificationv2

const (
	// The list of notification method.
	SMSNotifyMethod    Method = "sms"
	EmailNotifyMethod  Method = "email"
	VoiceNotifyMethod  Method = "voice"
	MobileNotifyMethod Method = "mobile"
)

// Method is a method of notification.
type Method string
