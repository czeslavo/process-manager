package manager

import "github.com/czeslavo/process-manager/1_voiding/messages"

type DocumentVoidingProcess struct {
	ID         string
	DocumentID string
	CustomerID string
	State      State
}

func NewDocumentVoidingProcess(processID, documentID, customerID string) DocumentVoidingProcess {
	return DocumentVoidingProcess{
		ID:         processID,
		DocumentID: documentID,
		CustomerID: customerID,
		State:      VoidingRequested,
	}
}

func (p *DocumentVoidingProcess) MarkingDocumentAsVoidedFailed() error {
	if err := p.State.canTransition(MarkingDocumentAsVoidedFailed); err != nil {
		return err
	}

	p.State = MarkingDocumentAsVoidedFailed
	return nil
}

func (p *DocumentVoidingProcess) MarkingDocumentAsVoidedSucceeded() error {
	if err := p.State.canTransition(MarkingDocumentAsVoidedSucceeded); err != nil {
		return err
	}

	p.State = MarkingDocumentAsVoidedSucceeded
	return nil
}

func (p *DocumentVoidingProcess) DocumentVoided() error {
	if err := p.State.canTransition(DocumentVoided); err != nil {
		return err
	}

	p.State = DocumentVoided
	return nil
}

func (p *DocumentVoidingProcess) AcknowledgeFailure() error {
	if err := p.State.canTransition(FailureAcknowledged); err != nil {
		return err
	}

	p.State = FailureAcknowledged
	return nil
}

func (p DocumentVoidingProcess) IsOngoing() bool {
	isTerminated := p.State == DocumentVoided || p.State == FailureAcknowledged
	return !isTerminated
}

func (p DocumentVoidingProcess) NextCommand() interface{} {
	switch p.State {
	case VoidingRequested:
		return &messages.MarkDocumentAsVoided{
			DocumentID:    p.DocumentID,
			RecipientID:   p.CustomerID,
			CorrelationID: p.ID,
		}
	case MarkingDocumentAsVoidedSucceeded:
		return &messages.VoidDocument{
			DocumentID:    p.DocumentID,
			RecipientID:   p.CustomerID,
			CorrelationID: p.ID,
		}
	default:
		return nil
	}
}
