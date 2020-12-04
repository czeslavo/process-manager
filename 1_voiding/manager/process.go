package manager

type DocumentVoidingProcess struct {
	ID         string
	DocumentID string
	State      State
}

func NewDocumentVoidingProcess(processID, documentID string) DocumentVoidingProcess {
	return DocumentVoidingProcess{
		ID:         processID,
		DocumentID: documentID,
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

func (p DocumentVoidingProcess) IsOngoing() bool {
	isTerminated := p.State == DocumentVoided || p.State == MarkingDocumentAsVoidedFailed
	return !isTerminated
}
