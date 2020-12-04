package messages

// Billing
type IssueDocument struct {
	DocumentID  string
	RecipientID string
	TotalAmount float64
}

type RequestDocumentVoiding struct {
	DocumentID string
}

type VoidDocument struct {
	DocumentID  string
	RecipientID string

	CorrelationID string
}

// Reports
type MarkDocumentAsVoided struct {
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
