package notificationv2

// Step is an action that will be added to notification rule.
type Step struct {
	ID        string     `json:"id"`
	Contact   Contact    `json:"contact,omitempty"`
	SendAfter SendAfter `json:"sendAfter,omitempty"`
	Enabled   bool       `json:"enabled,omitempty"`
}

// Contact defines the contact that notification will be sent to.
type Contact struct {
	Method Method `json:"method,omitempty"`
	To     string `json:"to,omitempty"`
}

// SendAfter defines minute time period notification will be sent after.
type SendAfter struct {
	TimeAmount int   `json:"timeAmount,omitempty"`
	TimeUnit   TimeUnit `json:"timeUnit,omitempty"`
}
