package billing

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/czeslavo/process-manager/1_voiding/messages"
	"github.com/pkg/errors"
)

type Service struct {
	repo     *DocumentsRepo
	eventBus *cqrs.EventBus
}

func NewService() *Service {
	return &Service{
		repo: NewDocumentsRepo(),
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
		messages.CommandHandlerFunc(
			"RequestDocumentVoiding",
			&messages.RequestDocumentVoiding{},
			s.handleRequestDocumentVoiding,
		),
		messages.CommandHandlerFunc(
			"VoidDocument",
			&messages.VoidDocument{},
			s.handleVoidDocument,
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

func (s *Service) handleRequestDocumentVoiding(ctx context.Context, cmd interface{}) error {
	requestDocumentVoiding := cmd.(*messages.RequestDocumentVoiding)
	fmt.Println("billing: request document voiding command")

	// Check if a document exists and is not voided yet. If it is, trigger the process of voiding.
	document, err := s.repo.GetByID(requestDocumentVoiding.DocumentID)
	if err != nil {
		return errors.Wrapf(err, "could not get document by id: %s", requestDocumentVoiding.DocumentID)
	}

	if document.IsVoided {
		// noop
		return nil
	}

	return s.eventBus.Publish(ctx, &messages.DocumentVoidRequested{
		DocumentID:  document.ID,
		RecipientID: document.RecipientID,
	})
}

func (s *Service) handleVoidDocument(ctx context.Context, cmd interface{}) error {
	requestDocumentVoiding := cmd.(*messages.VoidDocument)
	fmt.Println("billing: void document command")

	// Check if a document exists and is not voided yet. If it is, trigger the process of voiding.
	document, err := s.repo.GetByID(requestDocumentVoiding.DocumentID)
	if err != nil {
		return errors.Wrapf(err, "could not get document by id: %s", requestDocumentVoiding.DocumentID)
	}

	document.IsVoided = true
	s.repo.Store(document.ID, document)

	return s.eventBus.Publish(ctx, &messages.DocumentVoided{
		DocumentID:    document.ID,
		CorrelationID: requestDocumentVoiding.CorrelationID,
	})
}
