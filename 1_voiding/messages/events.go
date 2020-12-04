package messages

type DocumentIssued struct {
	CustomerID          string
	DocumentID          string
	DocumentTotalAmount float64
}

type DocumentVoidRequested struct {
	DocumentID  string
	RecipientID string
}

type MarkingDocumentAsVoidedSucceeded struct {
	DocumentID  string
	RecipientID string
}

type MarkingDocumentAsVoidedFailed struct {
	DocumentID  string
	RecipientID string
}
