package balance_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/czeslavo/process-manager/2_temporal/activities"
	"github.com/czeslavo/process-manager/2_temporal/balance"
	"github.com/czeslavo/process-manager/2_temporal/events"
	"github.com/czeslavo/process-manager/2_temporal/worker"
	"github.com/czeslavo/process-manager/2_temporal/workflows"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestTriggerReprocessTrip(t *testing.T) {
	ctx := context.Background()

	rand.Seed(time.Now().UnixNano())

	logger := watermill.NewStdLogger(true, true)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)
	reprocessingFinishedEvents, err := pubsub.Subscribe(ctx, events.TripReprocessingFinishedTopic)
	require.NoError(t, err)

	repo := balance.NewTripBalanceRepository()

	handler, err := balance.NewReprocessTripHandler(repo, pubsub)
	require.NoError(t, err)

	router := setupRouter(t, handler, pubsub, logger)
	go router.Run(ctx)
	<-router.Running()

	w, err := worker.NewWorker()
	require.NoError(t, err)
	require.NoError(t, workflows.RegisterWorkflows(w))
	require.NoError(t, activities.RegisterActivities(w, pubsub, repo))
	require.NoError(t, w.Start())

	balance := balance.NewTripBalance("trip-id", -0.54)
	repo.SaveBalance(balance)

	correlationID := uuid.NewString()
	err = handler.HandleReprocess(ctx, balance.TripUUID(), correlationID)
	require.NoError(t, err)

	<-reprocessingFinishedEvents
}

func setupRouter(
	t *testing.T,
	handler *balance.ReprocessTripHandler,
	pubsub *gochannel.GoChannel,
	logger watermill.LoggerAdapter,
) *message.Router {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	require.NoError(t, err)

	router.AddNoPublisherHandler("CreditNoteIssuedHandler", events.CreditNoteIssuedTopic, pubsub, handler.HandleCreditNoteIssued)
	router.AddNoPublisherHandler("RefundedHandler", events.RefundedTopic, pubsub, handler.HandleRefunded)
	return router
}
