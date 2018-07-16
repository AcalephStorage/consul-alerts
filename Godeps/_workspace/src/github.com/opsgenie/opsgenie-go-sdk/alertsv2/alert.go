package alertsv2

import (
	"time"
)

type Alert struct {
	ID             string      `json:"id,omitempty"`
	TinyID         string      `json:"tinyId,omitempty"`
	Alias          string      `json:"alias,omitempty"`
	Message        string      `json:"message,omitempty"`
	Status         string      `json:"status,omitempty"`
	Acknowledged   bool        `json:"acknowledged,omitempty"`
	IsSeen         bool        `json:"isSeen,omitempty"`
	Tags           []string    `json:"tags,omitempty"`
	Snoozed        bool        `json:"snoozed,omitempty"`
	SnoozedUntil   time.Time   `json:"snoozedUntil,omitempty"`
	Count          int         `json:"count,omitempty"`
	LastOccurredAt time.Time   `json:"lastOccuredAt,omitempty"`
	CreatedAt      time.Time   `json:"createdAt,omitempty"`
	UpdatedAt      time.Time   `json:"updatedAt,omitempty"`
	Source         string      `json:"source,omitempty"`
	Owner          string      `json:"owner,omitempty"`
	Priority       Priority    `json:"priority,omitempty"`
	Teams          []TeamMeta  `json:"teams,omitempty"`
	Integration    Integration `json:"integration,omitempty"`
	Report         Report      `json:"report,omitempty"`
}

type ListAlertResponse struct {
	ResponseMeta
	Alerts []Alert `json:"data"`
}
