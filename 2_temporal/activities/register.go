package activities

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
)

func RegisterActivities(w worker.Worker, pub message.Publisher) error {
	issueCreditNoteHandler := IssueCreditNoteHandler{publisher: pub}
	w.RegisterActivityWithOptions(issueCreditNoteHandler.IssueCreditNote, activity.RegisterOptions{
		Name: IssueCreditNoteActivity,
	})
	return nil
}
