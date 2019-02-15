package notificationv2

const (
	// Types of action that notification rule will have after creating. It is used to create an alert.
	CreateAlertActionType         ActionType = "create-alert"
	AcknowledgedAlertActionType   ActionType = "acknowledged-alert"
	ClosedAlertActionType         ActionType = "closed-alert"
	AssignedAlertActionType       ActionType = "assigned-alert"
	AddNoteActionType             ActionType = "add-note"
	ScheduleStartActionType       ActionType = "schedule-start"
	ScheduleEndActionType         ActionType = "schedule-end"
	IncomingCallRoutingActionType ActionType = "incoming-call-routing"
)

// ActionType is the type of notification action. Instead of
type ActionType string
