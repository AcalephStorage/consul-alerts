package userv2

// Escalation is a struct of escalation.
type Escalation struct {
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	OwnerTeam   OwnerTeam `json:"ownerTeam,omitempty"`
	Rules       []Rule    `json:"rules,omitempty"`
}

// OwnerTeam contains info about owner team of escalation.
type OwnerTeam struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Rule is a rule of escalation.
type Rule struct {
	Condition  string    `json:"condition,omitempty"`
	NotifyType string    `json:"notifyType,omitempty"`
	Delay      Delay     `json:"delay,omitempty"`
	Recipient  Recipient `json:"recipient,omitempty"`
}

// Delay contains info about delaying alerts.
type Delay struct {
	TimeAmount int    `json:"timeAmount,omitempty"`
	TimeUnit   string `json:"timeUnit,omitempty"`
}

// Recipient contains info about recipient of alerts.
type Recipient struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}
