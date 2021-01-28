package reports

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/czeslavo/process-manager/2_temporal/messages"
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

func (s Service) MarkDocumentAsVoided(ctx context.Context, documentUUID string) error {
	fmt.Println("reports: mark document as voided command")

	reports := s.repo.GetForDocument(documentUUID)
	for _, report := range reports {
		if report.IsPublished {
			return errors.New("cannot mark document as voided as report is already published")
		}
	}

	return nil
}

func (s *Service) handlePublishReport(ctx context.Context, cmd interface{}) error {
	publishReport := cmd.(*messages.PublishReport)
	fmt.Println("reports: publish report command")

	report := s.repo.GetOrCreate(publishReport.CustomerID)
	report.Publish()
	s.repo.Store(report)

	return nil
}
