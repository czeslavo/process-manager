package activities

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/czeslavo/process-manager/2_temporal/events"
	"github.com/google/uuid"
)

const (
	TriggerRefundActivity = "TriggerRefund"
)

func (h Handler) TriggerRefund(ctx context.Context, correlationID string) error {

	// triggering refund...

	e := events.Refunded{
		CorrelationID: correlationID,
	}

	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	err = h.publisher.Publish(events.RefundedTopic, message.NewMessage(
		uuid.NewString(),
		b,
	))

	return err
}
