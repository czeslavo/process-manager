package process

import (
	"fmt"

	"github.com/czeslavo/process-manager/1_voiding/messages"
	processmanager "github.com/czeslavo/process-manager/2_lib"
)

type DocumentVoidingProcess struct {
	id         string
	documentID string
	customerID string
	state      State
}

func NewDocumentVoidingProcess(processID string) *DocumentVoidingProcess {
	return &DocumentVoidingProcess{
		id: processID,
	}
}

func (p *DocumentVoidingProcess) Process(event processmanager.Event) error {
	switch t := event.Payload.(type) {
	case DocumentVoidRequested:
		return p.setState(VoidingRequestedState)
	case MarkingDocumentAsVoidedFailed:
		return p.setState(MarkingDocumentAsVoidedFailedState)
	case MarkingDocumentAsVoidedSucceeded:
		return p.setState(MarkingDocumentAsVoidedSucceededState)
	case DocumentVoided:
		return p.setState(DocumentVoidedState)
	default:
		return fmt.Errorf("unknown event: %T", t)
	}
}

func (p *DocumentVoidingProcess) setState(state State) error {
	if err := p.state.canTransition(state); err != nil {
		return err
	}
	p.state = state
	return nil
}

func (p *DocumentVoidingProcess) NextCommands() []interface{} {
	panic("implement me")
}

func (p *DocumentVoidingProcess) ID() string {
	panic("implement me")
}

func (p *DocumentVoidingProcess) State() string {
	panic("implement me")
}

func (p DocumentVoidingProcess) nextCommand() interface{} {
	switch p.state {
	case VoidingRequestedState:
		return &messages.MarkDocumentAsVoided{
			DocumentID:    p.documentID,
			RecipientID:   p.customerID,
			CorrelationID: p.id,
		}
	case MarkingDocumentAsVoidedSucceededState:
		return &messages.VoidDocument{
			DocumentID:    p.documentID,
			RecipientID:   p.customerID,
			CorrelationID: p.id,
		}
	default:
		return nil
	}
}
