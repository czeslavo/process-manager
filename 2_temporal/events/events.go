package events

const (
	CreditNoteIssuedTopic         = "credit-note-issued"
	TripReprocessingFinishedTopic = "trip-reprocessing-finished"
	RefundedTopic                 = "refunded"
)

type CreditNoteIssued struct {
	CorrelationID string `json:"correlation_id"`
}

type TripReprocessingFinished struct {
	CorrelationID string `json:"correlation_id"`
}

type Refunded struct {
	CorrelationID string `json:"correlation_id"`
}
