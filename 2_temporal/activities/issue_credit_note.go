package activities

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/czeslavo/process-manager/2_temporal/events"

	"github.com/ThreeDotsLabs/watermill/message"
)

const (
	IssueCreditNoteActivity = "IssueCreditNote"
)

type IssueCreditNoteHandler struct {
	publisher message.Publisher
}

func (h IssueCreditNoteHandler) IssueCreditNote(ctx context.Context, correlationID string) error {
	e := events.CreditNoteIssued{
		CorrelationID: correlationID,
	}

	// issuing credit note...

	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	err = h.publisher.Publish(events.EventsTopic, message.NewMessage(
		uuid.NewString(),
		b,
	))

	return err
}
