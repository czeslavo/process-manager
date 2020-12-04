package manager

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/czeslavo/process-manager/1_voiding/messages"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type DocumentVoidingProcessManager struct {
	commandBus *cqrs.CommandBus
	repo       *Repo
	slowedDown bool
}

func NewDocumentVoidingProcessManager() *DocumentVoidingProcessManager {
	return &DocumentVoidingProcessManager{
		repo:       NewRepo(),
		slowedDown: os.Getenv("SLOW_DOWN") == "1",
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

func (m DocumentVoidingProcessManager) CommandHandlers() []cqrs.CommandHandler {
	return []cqrs.CommandHandler{
		messages.CommandHandlerFunc(
			"AcknowledgeProcessFailure",
			&messages.AcknowledgeProcessFailure{},
			m.acknowledgeProcessFailure,
		),
	}
}

func (m DocumentVoidingProcessManager) GetProcessForDocument(documentID string) (DocumentVoidingProcess, bool) {
	return m.repo.GetOngoingForDocument(documentID)
}

func (m DocumentVoidingProcessManager) GetAllOngoingOrFailed() []DocumentVoidingProcess {
	return m.repo.GetAllOngoing()
}

func (m DocumentVoidingProcessManager) handleDocumentVoidRequested(ctx context.Context, event interface{}) error {
	documentVoidRequested := event.(*messages.DocumentVoidRequested)
	fmt.Println("manager: document void requested event")

	if m.repo.IsOngoingForDocument(documentVoidRequested.DocumentID) {
		// only one process can be in progress for a single document
		return nil
	}

	processID := uuid.New().String()
	process := NewDocumentVoidingProcess(processID, documentVoidRequested.DocumentID, documentVoidRequested.RecipientID)

	m.slowDownIfConfigured()
	if cmd := process.NextCommand(); cmd != nil {
		if err := m.commandBus.Send(ctx, cmd); err != nil {
			return errors.Wrap(err, "failed to send command")
		}
	}

	m.repo.Store(process)

	return nil
}

func (m DocumentVoidingProcessManager) handleMarkingDocumentAsVoidedSucceeded(ctx context.Context, event interface{}) error {
	markingDocumentAsVoidedSucceeded := event.(*messages.MarkingDocumentAsVoidedSucceeded)
	fmt.Println("manager: marking document as voided succeeded event")

	process, err := m.repo.GetProcess(markingDocumentAsVoidedSucceeded.CorrelationID)
	if err != nil {
		return err
	}
	if err := process.MarkingDocumentAsVoidedSucceeded(); err != nil {
		return err
	}

	m.slowDownIfConfigured()
	if cmd := process.NextCommand(); cmd != nil {
		if err := m.commandBus.Send(ctx, cmd); err != nil {
			return errors.Wrap(err, "failed to send command")
		}
	}

	m.repo.Store(process)

	return nil
}

func (m DocumentVoidingProcessManager) handleMarkingDocumentAsVoidedFailed(ctx context.Context, event interface{}) error {
	markingDocumentAsVoidedFailed := event.(*messages.MarkingDocumentAsVoidedFailed)
	fmt.Println("manager: marking document as voided failed event")

	process, err := m.repo.GetProcess(markingDocumentAsVoidedFailed.CorrelationID)
	if err != nil {
		return err
	}
	if err := process.MarkingDocumentAsVoidedFailed(); err != nil {
		return err
	}

	m.slowDownIfConfigured()
	m.repo.Store(process)

	return nil
}

func (m DocumentVoidingProcessManager) handleDocumentVoided(ctx context.Context, event interface{}) error {
	documentVoided := event.(*messages.DocumentVoided)
	fmt.Println("manager: document voided event")

	process, err := m.repo.GetProcess(documentVoided.CorrelationID)
	if err != nil {
		return err
	}
	if err := process.DocumentVoided(); err != nil {
		return err
	}

	m.slowDownIfConfigured()
	m.repo.Store(process)

	return nil
}

func (m DocumentVoidingProcessManager) acknowledgeProcessFailure(_ context.Context, cmd interface{}) error {
	acknowledgeFailure := cmd.(*messages.AcknowledgeProcessFailure)
	fmt.Println("manager: document voided event")

	process, err := m.repo.GetProcess(acknowledgeFailure.ProcessID)
	if err != nil {
		return err
	}
	if err := process.AcknowledgeFailure(); err != nil {
		return err
	}

	m.slowDownIfConfigured()
	m.repo.Store(process)

	return nil
}

func (m DocumentVoidingProcessManager) slowDownIfConfigured() {
	if m.slowedDown {
		time.Sleep(time.Second * 5)
	}
}
