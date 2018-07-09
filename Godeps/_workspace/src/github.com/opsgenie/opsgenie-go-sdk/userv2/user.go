package userv2

import "time"

// User contains user data.
type User struct {
	ID           string              `json:"id,omitempty"`
	Blocked      bool                `json:"blocked,omitempty"`
	Verified     bool                `json:"verified,omitempty"`
	Username     string              `json:"username,omitempty"`
	FullName     string              `json:"fullName,omitempty"`
	Role         UserRole            `json:"role,omitempty"`
	TimeZone     string              `json:"timeZone,omitempty"`
	Locale       string              `json:"locale,omitempty"`
	UserAddress  UserAddress         `json:"userAddress,omitempty"`
	CreatedAt    time.Time           `json:"createdAt,omitempty"`
	MutedUntil   time.Time           `json:"mutedUntil,omitempty"`
	Details      map[string][]string `json:"details,omitempty"`
	Tags         []string            `json:"tags,omitempty"`
	UserContacts []UserContact       `json:"userContacts,omitempty"`
}
