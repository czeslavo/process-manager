package manager

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/czeslavo/process-manager/1_voiding/messages"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type DocumentVoidingProcessManager struct {
	commandBus *cqrs.CommandBus
	repo       *Repo
}

func NewDocumentVoidingProcessManager() *DocumentVoidingProcessManager {
	return &DocumentVoidingProcessManager{
		repo: NewRepo(),
	}
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
		messages.EventHandlerFunc(
			"DocumentVoided",
			&messages.DocumentVoided{},
			m.handleDocumentVoided,
		),
	}
}

func (m DocumentVoidingProcessManager) GetProcessForDocument(documentID string) (DocumentVoidingProcess, bool) {
	return m.repo.GetOngoingForDocument(documentID)
}

func (m DocumentVoidingProcessManager) handleDocumentVoidRequested(ctx context.Context, event interface{}) error {
	documentVoidRequested := event.(*messages.DocumentVoidRequested)
	fmt.Println("manager: document void requested event")

	if m.repo.IsOngoingForDocument(documentVoidRequested.DocumentID) {
		// only one process can be in progress for a single document
		return nil
	}

	processID := uuid.New().String()
	process := m.repo.GetOrCreateProcess(processID, documentVoidRequested.DocumentID)

	if err := m.commandBus.Send(ctx, &messages.MarkDocumentAsVoided{
		DocumentID:  documentVoidRequested.DocumentID,
		RecipientID: documentVoidRequested.RecipientID,
	}); err != nil {
		return errors.Wrap(err, "failed to send command")
	}

	m.repo.Store(process)

	return nil
}

func (m DocumentVoidingProcessManager) handleMarkingDocumentAsVoidedSucceeded(ctx context.Context, event interface{}) error {
	markingDocumentAsVoidedSucceeded := event.(*messages.MarkingDocumentAsVoidedSucceeded)
	fmt.Println("manager: marking document as voided succeeded event")

	processID := markingDocumentAsVoidedSucceeded.CorrelationID
	process := m.repo.GetOrCreateProcess(
		processID,
		markingDocumentAsVoidedSucceeded.DocumentID,
	)
	if err := process.MarkingDocumentAsVoidedSucceeded(); err != nil {
		return err
	}

	if err := m.commandBus.Send(ctx, &messages.VoidDocument{
		DocumentID:  markingDocumentAsVoidedSucceeded.DocumentID,
		RecipientID: markingDocumentAsVoidedSucceeded.RecipientID,
	}); err != nil {
		return errors.Wrap(err, "failed to send command")
	}

	m.repo.Store(process)

	return nil
}

func (m DocumentVoidingProcessManager) handleMarkingDocumentAsVoidedFailed(ctx context.Context, event interface{}) error {
	markingDocumentAsVoidedFailed := event.(*messages.MarkingDocumentAsVoidedFailed)
	fmt.Println("manager: marking document as voided failed event")

	processID := markingDocumentAsVoidedFailed.CorrelationID
	process := m.repo.GetOrCreateProcess(
		processID,
		markingDocumentAsVoidedFailed.DocumentID,
	)
	if err := process.MarkingDocumentAsVoidedFailed(); err != nil {
		return err
	}

	m.repo.Store(process)

	return nil
}

func (m DocumentVoidingProcessManager) handleDocumentVoided(ctx context.Context, event interface{}) error {
	documentVoided := event.(*messages.DocumentVoided)
	fmt.Println("manager: document voided event")

	processID := documentVoided.CorrelationID
	process := m.repo.GetOrCreateProcess(
		processID,
		documentVoided.DocumentID,
	)
	if err := process.MarkingDocumentAsVoidedFailed(); err != nil {
		return err
	}

	m.repo.Store(process)

	return nil
}
