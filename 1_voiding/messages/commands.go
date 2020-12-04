package messages

// Billing
type IssueDocument struct {
	DocumentID  string
	RecipientID string
	TotalAmount float64
}

type CompleteDocumentVoiding struct {
	DocumentID  string
	RecipientID string
}

type AbortDocumentVoiding struct {
	DocumentID  string
	RecipientID string
}

// Reports
type MarkDocumentAsVoided struct {
	DocumentID  string
	RecipientID string
}
