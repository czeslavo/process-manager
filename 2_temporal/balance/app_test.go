package balance

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTriggerReprocessTrip(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	repo := NewTripBalanceRepository()

	balance := NewTripBalance("trip-id", -0.54)
	repo.SaveBalance(balance)

	handler := &TriggerReprocessTripHandler{
		repo: repo,
	}
	err := handler.Handle(ctx, balance.TripUUID(), "correlation-id")
	require.NoError(t, err)
}
