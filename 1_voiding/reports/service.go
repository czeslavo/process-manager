package reports

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/czeslavo/process-manager/1_voiding/messages"
)

type Service struct {
	repo *Repo
}

func NewService() *Service {
	return &Service{
		repo: NewRepo(),
	}
}

func (s Service) CommandHandlers() []cqrs.CommandHandler {
	return nil
}

func (s Service) EventHandlers() []cqrs.EventHandler {
	return []cqrs.EventHandler{
		messages.EventHandlerFunc(
			"DocumentIssued",
			&messages.DocumentIssued{},
			s.handleDocumentIssued,
		),
	}
}

func (s Service) handleDocumentIssued(_ context.Context, event interface{}) error {
	documentIssued := event.(*messages.DocumentIssued)
	fmt.Println("reports: document issued event")

	s.repo.AppendToReport(documentIssued.CustomerID, documentIssued.DocumentID, documentIssued.DocumentTotalAmount)

	return nil
}

func (s Service) handleMarkDocumentAsVoided(_ context.Context, cmd interface{}) error {
	markDocumentAsVoided := cmd.(*messages.MarkDocumentAsVoided)
	_ = markDocumentAsVoided
	return nil
}
