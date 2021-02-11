package balance

import (
	"context"

	"github.com/czeslavo/process-manager/2_temporal/workflows"
	"github.com/sirupsen/logrus"
	"go.temporal.io/sdk/client"
)

const taskQueueName = "trip_balance"

type temporal struct {
	client client.Client
	logger logrus.FieldLogger
}

func (c temporal) startWorkflow(ctx context.Context, correlationID string, reprocessType ReprocessType) error {
	var workflowName string
	switch reprocessType {
	case ReprocessTypeRefund:
		workflowName = workflows.RefundWorkflowName
	}

	run, err := c.client.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:        correlationID,
			TaskQueue: taskQueueName,
		},
		workflowName,
	)
	if err != nil {
		return err
	}

	c.logger.WithField("correlation_id", correlationID).WithField("run_id", run.GetRunID()).Info("Run started")
	return nil
}

type TriggerReprocessTripHandler struct {
	repo *TripBalanceRepository
}

func (h TriggerReprocessTripHandler) Handle(ctx context.Context, tripUUID, correlationID string) error {
	logger := logrus.New()
	c, err := client.NewClient(client.Options{})
	if err != nil {
		return err
	}

	balance, err := h.repo.GetBalance(tripUUID)
	if err != nil {
		return err
	}

	reprocessType, err := balance.ReprocessTrip(correlationID)
	if err != nil {
		return err
	}

	t := temporal{
		client: c,
		logger: logger,
	}

	return t.startWorkflow(ctx, correlationID, reprocessType)
}
