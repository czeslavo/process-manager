package activities

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/czeslavo/process-manager/2_temporal/events"
	"github.com/google/uuid"
)

const (
	ReprocessingFinishedActivity = "ReprocessingFinished"
)

func (h Handler) ReprocessingFinished(ctx context.Context, tripUUID, amendmentID string) error {
	balance, err := h.repo.GetBalance(tripUUID)
	if err != nil {
		return err
	}

	if err := balance.ReprocessFinished(amendmentID); err != nil {
		return err
	}

	b, err := json.Marshal(events.TripReprocessingFinished{CorrelationID: amendmentID})
	if err != nil {
		return err
	}
	if err := h.publisher.Publish(events.TripReprocessingFinishedTopic, message.NewMessage(uuid.NewString(), b)); err != nil {
		return err
	}

	return nil
}
