package userv2

import "time"

// ListUserForwardingRulesResponse is a response with list of forwarding rules.
type ListUserForwardingRulesResponse struct {
	ForwardingRules []ForwardingRule `json:"data,omitempty"`
	ResponseMeta
}

// ForwardingRule contains data of forwarding rule.
type ForwardingRule struct {
	ID        string    `json:"id,omitempty"`
	Alias     string    `json:"alias,omitempty"`
	FromUser  UserMeta  `json:"fromUser,omitempty"`
	ToUser    UserMeta  `json:"toUser,omitempty"`
	StartDate time.Time `json:"startDate,omitempty"`
	EndDate   time.Time `json:"endDate,omitempty"`
}
