package billing

import (
	"context"
	"fmt"

	"github.com/czeslavo/process-manager/2_temporal/voiding"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/czeslavo/process-manager/2_temporal/messages"
	"github.com/pkg/errors"
	"go.temporal.io/sdk/client"
)

type Service struct {
	repo           *DocumentsRepo
	eventBus       *cqrs.EventBus
	temporalClient client.Client
}

func NewService(temporalClient client.Client) *Service {
	return &Service{
		repo:           NewDocumentsRepo(),
		temporalClient: temporalClient,
	}
}

func (s *Service) CommandHandlers(eventBus *cqrs.EventBus) []cqrs.CommandHandler {
	s.eventBus = eventBus

	return []cqrs.CommandHandler{
		messages.CommandHandlerFunc(
			"IssueDocument",
			&messages.IssueDocument{},
			s.handleIssueDocument,
		),
	}
}

func (s Service) EventHandlers() []cqrs.EventHandler {
	return nil
}

func (s Service) GetDocument(id string) (Document, error) {
	return s.repo.GetByID(id)
}

func (s Service) handleIssueDocument(ctx context.Context, cmd interface{}) error {
	issueDocumentCmd := cmd.(*messages.IssueDocument)
	fmt.Println("billing: issue document command")

	if _, err := s.repo.GetByID(issueDocumentCmd.DocumentID); err == nil {
		// document already issued
		return nil
	}

	document := Document{
		ID:          issueDocumentCmd.DocumentID,
		RecipientID: issueDocumentCmd.RecipientID,
		TotalAmount: issueDocumentCmd.TotalAmount,
	}
	s.repo.Store(issueDocumentCmd.DocumentID, document)

	if err := s.eventBus.Publish(ctx, &messages.DocumentIssued{
		CustomerID:          document.RecipientID,
		DocumentID:          document.ID,
		DocumentTotalAmount: document.TotalAmount,
	}); err != nil {
		return errors.Wrap(err, "failed to publish event")
	}

	return nil
}

func (s *Service) TriggerDocumentVoidingWorkflow(ctx context.Context, documentID string) error {
	fmt.Println("billing: request document voiding command")

	// Check if a document exists and is not voided yet. If it is, trigger the process of voiding.
	document, err := s.repo.GetByID(documentID)
	if err != nil {
		return errors.Wrapf(err, "could not get document by id: %s", documentID)
	}

	if document.IsVoided {
		// noop
		return nil
	}

	if _, err = s.temporalClient.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{TaskQueue: voiding.VoidDocumentsWorkflowQueueName},
		voiding.VoidDocumentsWorkflowName,
		documentID,
	); err != nil {
		return errors.Wrap(err, "could not start voiding workflow in temporal")
	}

	return nil
}

func (s *Service) VoidDocument(ctx context.Context, documentID string) error {
	fmt.Println("billing: void document command")

	document, err := s.repo.GetByID(documentID)
	if err != nil {
		return errors.Wrapf(err, "could not get document by id: %s", documentID)
	}

	document.IsVoided = true
	s.repo.Store(document.ID, document)

	return nil
}
