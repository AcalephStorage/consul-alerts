package notificationv2

const (
	// The list of fields, which are used for build filter of alerts.
	ActionsField         Field = "actions"
	AliasField           Field = "alias"
	DescriptionField     Field = "description"
	EntityField          Field = "entity"
	MessageField         Field = "message"
	RecipientsField      Field = "recipients"
	SourceField          Field = "source"
	TeamsField           Field = "teams"
	ExtraPropertiesField Field = "extra-properties"
)

// Field is the name of alert field, which is used to build filter of alerts.
type Field string
