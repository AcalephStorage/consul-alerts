package notificationv2

const (
	// The list of matches which are used for creating criteria of notification.
	MatchAllType           ConditionType = "match-all"
	MatchAnyConditionsType ConditionType = "match-any-condition"
	MatchAllConditionsType ConditionType = "match-all-conditions"
)

// ConditionType is a type of matching of alert fields, which is used for building criteria of notification.
type ConditionType string
