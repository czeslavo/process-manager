package messages

// Billing
type IssueDocument struct {
	DocumentID  string
	RecipientID string
	TotalAmount float64
}

type VoidDocument struct {
	DocumentID  string
	RecipientID string

	CorrelationID string
}

type PublishReport struct {
	CustomerID string
}

// Process manager
type AcknowledgeProcessFailure struct {
	ProcessID string
}
