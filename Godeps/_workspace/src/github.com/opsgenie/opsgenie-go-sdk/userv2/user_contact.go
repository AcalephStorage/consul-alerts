package userv2

// UserContact contains data of user contact.
type UserContact struct {
	ID             string        `json:"id"`
	To             string        `json:"to"`
	ContactMethod  ContactMethod `json:"contactMethod"`
	DisabledReason string        `json:"disabledReason"`
	Enabled        bool          `json:"enabled"`
}
