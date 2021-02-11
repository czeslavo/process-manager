package events

const (
	EventsTopic = "events"
)

type CreditNoteIssued struct {
	CorrelationID string `json:"correlation_id"`
}
