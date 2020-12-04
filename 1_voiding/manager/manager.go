package manager

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/czeslavo/process-manager/1_voiding/messages"
)

type DocumentVoidingProcessManager struct {
	commandBus *cqrs.CommandBus
}

func NewDocumentVoidingProcessManager() *DocumentVoidingProcessManager {
	return &DocumentVoidingProcessManager{}
}

func (m DocumentVoidingProcessManager) EventHandlers(commandBus *cqrs.CommandBus) []cqrs.EventHandler {
	m.commandBus = commandBus

	return []cqrs.EventHandler{
		messages.EventHandlerFunc(
			"DocumentVoidRequested",
			&messages.DocumentVoidRequested{},
			m.handleDocumentVoidRequested,
		),
		messages.EventHandlerFunc(
			"MarkingDocumentAsVoidedSucceeded",
			&messages.MarkingDocumentAsVoidedSucceeded{},
			m.handleMarkingDocumentAsVoidedSucceeded,
		),
		messages.EventHandlerFunc(
			"MarkingDocumentAsVoidedFailed",
			&messages.MarkingDocumentAsVoidedFailed{},
			m.handleMarkingDocumentAsVoidedFailed,
		),
	}
}

func (m DocumentVoidingProcessManager) handleDocumentVoidRequested(ctx context.Context, event interface{}) error {
	documentVoidRequested := event.(*messages.DocumentVoidRequested)
	fmt.Println("manager: document void requested event")

	if err := m.commandBus.Send(ctx, &messages.MarkDocumentAsVoided{
		DocumentID:  documentVoidRequested.DocumentID,
		RecipientID: documentVoidRequested.RecipientID,
	}); err != nil {
		return errors.Wrap(err, "failed to send command")
	}

	return nil
}

func (m DocumentVoidingProcessManager) handleMarkingDocumentAsVoidedSucceeded(ctx context.Context, event interface{}) error {
	markingDocumentAsVoidedSucceeded := event.(*messages.MarkingDocumentAsVoidedSucceeded)
	fmt.Println("manager: marking document as voided succeeded event")

	if err := m.commandBus.Send(ctx, &messages.CompleteDocumentVoiding{
		DocumentID:  markingDocumentAsVoidedSucceeded.DocumentID,
		RecipientID: markingDocumentAsVoidedSucceeded.RecipientID,
	}); err != nil {
		return errors.Wrap(err, "failed to send command")
	}
	return nil
}

func (m DocumentVoidingProcessManager) handleMarkingDocumentAsVoidedFailed(ctx context.Context, event interface{}) error {
	markingDocumentAsVoidedFailed := event.(*messages.MarkingDocumentAsVoidedFailed)
	fmt.Println("manager: marking document as voided failed event")

	if err := m.commandBus.Send(ctx, &messages.AbortDocumentVoiding{
		DocumentID:  markingDocumentAsVoidedFailed.DocumentID,
		RecipientID: markingDocumentAsVoidedFailed.RecipientID,
	}); err != nil {
		return errors.Wrap(err, "failed to send command")
	}
	return nil
}
