package balance

import (
	"context"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/stretchr/testify/require"
)

func TestTriggerReprocessTrip(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	repo := NewTripBalanceRepository()

	balance := NewTripBalance("trip-id", -0.54)
	repo.SaveBalance(balance)

	logger := watermill.NewStdLogger(true, true)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)

	handler, err := NewReprocessTripHandler(repo, pubsub)
	require.NoError(t, err)

	router := setupRouter(t, handler, pubsub, logger)
	go router.Run(ctx)
	<-router.Running()

	err = handler.HandleReprocess(ctx, balance.TripUUID(), "correlation-id")
	require.NoError(t, err)
}

func setupRouter(
	t *testing.T,
	handler *ReprocessTripHandler,
	pubsub *gochannel.GoChannel,
	logger watermill.LoggerAdapter,
) *message.Router {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	require.NoError(t, err)

	router.AddNoPublisherHandler("CreditNoteIssuedHandler", "CreditNoteIssued", pubsub, handler.HandleCreditNoteIssued)
	return router
}
