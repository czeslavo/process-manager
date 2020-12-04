package reports

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/czeslavo/process-manager/1_voiding/messages"
)

type Service struct {
	repo     *Repo
	eventBus *cqrs.EventBus
}

func NewService() *Service {
	return &Service{
		repo: NewRepo(),
	}
}

func (s *Service) CommandHandlers(eventBus *cqrs.EventBus) []cqrs.CommandHandler {
	s.eventBus = eventBus

	return []cqrs.CommandHandler{
		messages.CommandHandlerFunc(
			"MarkDocumentAsVoided",
			&messages.MarkDocumentAsVoided{},
			s.handleMarkDocumentAsVoided,
		),
		messages.CommandHandlerFunc(
			"PublishReport",
			&messages.PublishReport{},
			s.handlePublishReport,
		),
	}
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

func (s Service) GetReports() []Report {
	var reports []Report
	for _, r := range s.repo.customerReport {
		reports = append(reports, r)
	}
	return reports
}

func (s Service) handleDocumentIssued(_ context.Context, event interface{}) error {
	documentIssued := event.(*messages.DocumentIssued)
	fmt.Println("reports: document issued event")

	report := s.repo.GetOrCreate(documentIssued.CustomerID)
	report.AppendDocument(documentIssued.DocumentID, documentIssued.DocumentTotalAmount)
	s.repo.Store(report)

	return nil
}

func (s Service) handleMarkDocumentAsVoided(ctx context.Context, cmd interface{}) error {
	markDocumentAsVoided := cmd.(*messages.MarkDocumentAsVoided)
	fmt.Println("reports: mark document as voided command")

	report := s.repo.GetOrCreate(markDocumentAsVoided.RecipientID)
	if report.IsPublished {
		return s.eventBus.Publish(ctx, &messages.MarkingDocumentAsVoidedFailed{
			DocumentID:  markDocumentAsVoided.DocumentID,
			RecipientID: markDocumentAsVoided.RecipientID,
		})
	}

	return s.eventBus.Publish(ctx, &messages.MarkingDocumentAsVoidedSucceeded{
		DocumentID:  markDocumentAsVoided.DocumentID,
		RecipientID: markDocumentAsVoided.RecipientID,
	})
}

func (s *Service) handlePublishReport(ctx context.Context, cmd interface{}) error {
	publishReport := cmd.(*messages.PublishReport)
	fmt.Println("reports: publish report command")

	report := s.repo.GetOrCreate(publishReport.CustomerID)
	report.Publish()
	s.repo.Store(report)

	return nil
}
