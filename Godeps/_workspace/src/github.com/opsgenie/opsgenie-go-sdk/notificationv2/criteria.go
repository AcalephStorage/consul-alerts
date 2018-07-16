package notificationv2

// Criteria defines the conditions that will be checked before applying notification rules and type of the operations
// that will be applied on these conditions.
type Criteria struct {
	Type       ConditionType `json:"type"`
	Conditions []Condition   `json:"conditions"`
}

// Condition defines the conditions that will be checked before applying notification rules.
type Condition struct {
	Field         Field     `json:"field"`
	Key           string    `json:"key"`
	Not           bool      `json:"not"`
	Operation     Operation `json:"operation"`
	ExpectedValue string    `json:"expectedValue"`
	Order         int       `json:"order"`
}
