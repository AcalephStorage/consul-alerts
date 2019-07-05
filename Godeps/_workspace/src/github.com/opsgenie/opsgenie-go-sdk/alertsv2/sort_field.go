package alertsv2

type SortField string

const (
	CreatedAt       SortField = "createdAt"
	UpdatedAt       SortField = "updatedAt"
	TinyId          SortField = "tinyId"
	Alias           SortField = "alias"
	Message         SortField = "message"
	Status          SortField = "status"
	Acknowledged    SortField = "acknowledged"
	IsSeen          SortField = "isSeen"
	Snoozed         SortField = "snoozed"
	SnoozedUntil    SortField = "snoozedUntil"
	Count           SortField = "count"
	LastOccurredAt  SortField = "lastOccuredAt"
	Source          SortField = "source"
	Owner           SortField = "owner"
	IntegrationName SortField = "integration.name"
	IntegrationType SortField = "integration.type"
	AckTime         SortField = "report.ackTime"
	CloseTime       SortField = "report.closeTime"
	AcknowledgedBy  SortField = "report.acknowledgedBy"
	ClosedBy        SortField = "report.closedBy"
)
