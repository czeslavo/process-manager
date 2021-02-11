package balance

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/czeslavo/process-manager/2_temporal/events"

	"github.com/google/uuid"

	"github.com/czeslavo/process-manager/2_temporal/activities"
	"github.com/czeslavo/process-manager/2_temporal/workflows"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/czeslavo/process-manager/2_temporal/worker"
	"github.com/stretchr/testify/require"
)

func TestTriggerReprocessTrip(t *testing.T) {
	ctx := context.Background()

	rand.Seed(time.Now().UnixNano())

	logger := watermill.NewStdLogger(true, true)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)

	repo := NewTripBalanceRepository()

	handler, err := NewReprocessTripHandler(repo, pubsub)
	require.NoError(t, err)

	router := setupRouter(t, handler, pubsub, logger)
	go router.Run(ctx)
	<-router.Running()

	w, err := worker.NewWorker()
	require.NoError(t, err)
	require.NoError(t, workflows.RegisterWorkflows(w))
	require.NoError(t, activities.RegisterActivities(w, pubsub))
	require.NoError(t, w.Start())

	balance := NewTripBalance("trip-id", -0.54)
	repo.SaveBalance(balance)

	correlationID := uuid.NewString()
	err = handler.HandleReprocess(ctx, balance.TripUUID(), correlationID)
	require.NoError(t, err)

	time.Sleep(10 * time.Minute)
}

func setupRouter(
	t *testing.T,
	handler *ReprocessTripHandler,
	pubsub *gochannel.GoChannel,
	logger watermill.LoggerAdapter,
) *message.Router {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	require.NoError(t, err)

	router.AddNoPublisherHandler("CreditNoteIssuedHandler", events.EventsTopic, pubsub, handler.HandleCreditNoteIssued)
	return router
}
