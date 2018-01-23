package notificationv2

const (
	// The list of condition operation, which are used for build notification criteria.
	MatchesConditionOperation      				Operation = "matches"
	EqualsConditionOperation  					Operation = "equals"
	IsEmptyConditionOperation 					Operation = "is-empty"
	ContainsConditionOperation 					Operation = "contains"
	StartsWithConditionOperation 				Operation = "starts-with"
	EqualsIgnoreWhitespaceConditionOperation 	Operation = "equals-ignore-whitespace"
	GreaterThanConditionOperation 				Operation = "greater-than"
	LessThanConditionOperation 					Operation = "less-than"
)

// Operation is type of matching operation, which is used to build filter of alerts.
type Operation string
