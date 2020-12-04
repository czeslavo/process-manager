package billing

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/czeslavo/process-manager/1_voiding/messages"
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
	}
}

func (s Service) EventHandlers() []cqrs.EventHandler {
	return nil
}

func (s Service) handleIssueDocument(ctx context.Context, cmd interface{}) error {
	issueDocumentCmd := cmd.(*messages.IssueDocument)
	fmt.Printf("billing: issue document command: %v\n", cmd)

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
