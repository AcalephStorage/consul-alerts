package userv2

const (
	// The list of notification methods.
	SMSContactMethod    ContactMethod = "sms"
	EmailContactMethod  ContactMethod = "email"
	VoiceContactMethod  ContactMethod = "voice"
	MobileContactMethod ContactMethod = "mobile"
)

// ContactMethod is a type of user contact method.
type ContactMethod string
