package messages

// Billing
type DocumentIssued struct {
	CustomerID          string
	DocumentID          string
	DocumentTotalAmount float64
}

type DocumentVoidRequested struct {
	DocumentID  string
	RecipientID string
}

type DocumentVoided struct {
	DocumentID string

	CorrelationID string
}

// Reports
type MarkingDocumentAsVoidedSucceeded struct {
	DocumentID  string
	RecipientID string

	CorrelationID string
}

type MarkingDocumentAsVoidedFailed struct {
	DocumentID  string
	RecipientID string

	CorrelationID string
}
