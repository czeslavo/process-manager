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

func (h Handler) IssueCreditNote(ctx context.Context, correlationID string) error {
	e := events.CreditNoteIssued{
		CorrelationID: correlationID,
	}

	// issuing credit note...

	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	err = h.publisher.Publish(events.CreditNoteIssuedTopic, message.NewMessage(
		uuid.NewString(),
		b,
	))

	return err
}
