package activities

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/czeslavo/process-manager/2_temporal/balance"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

func RegisterActivities(w worker.Worker, pub message.Publisher, repo *balance.TripBalanceRepository) error {
	handler := NewHandler(pub, repo)

	w.RegisterActivityWithOptions(handler.IssueCreditNote, activity.RegisterOptions{
		Name: IssueCreditNoteActivity,
	})

	w.RegisterActivityWithOptions(handler.TriggerRefund, activity.RegisterOptions{
		Name: TriggerRefundActivity,
	})

	w.RegisterActivityWithOptions(handler.ReprocessingFinished, activity.RegisterOptions{
		Name: ReprocessingFinishedActivity,
	})

	return nil
}
