package balance

import (
	"context"

	"github.com/czeslavo/process-manager/2_temporal/workflows"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.temporal.io/sdk/client"
)

const taskQueueName = "trip_balance"

type temporal struct {
	client client.Client
	logger logrus.FieldLogger
}

func (c temporal) startWorkflow(ctx context.Context, correlationID string, reprocessType ReprocessType) error {
	workflowName, err := getWorkflowName(reprocessType)

	run, err := c.client.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:        correlationID,
			TaskQueue: taskQueueName,
		},
		workflowName,
		// args?
	)
	if err != nil {
		return err
	}

	c.logger.WithField("correlation_id", correlationID).WithField("run_id", run.GetRunID()).Info("Run started")
	return nil
}

func getWorkflowName(reprocessType ReprocessType) (string, error) {
	switch reprocessType {
	case ReprocessTypeRefund:
		return workflows.RefundWorkflowName, nil
	case ReprocessTypeAdditionalPayment:
		return workflows.ChargeAdditionalAmountWorkflowName, nil
	case ReprocessTypeChangedContractDetails:
		return workflows.ReissueDocumentsWorkflowName, nil
	default:
		return "", errors.New("unsupported reprocess type")
	}
}

type ReprocessTripHandler struct {
	repo *TripBalanceRepository
}

func (h ReprocessTripHandler) HandleReprocess(ctx context.Context, tripUUID, correlationID string) error {
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

func (h ReprocessTripHandler) HandleReprocessFinished(ctx context.Context, tripUUID, correlationID string) error {
	balance, err := h.repo.GetBalance(tripUUID)
	if err != nil {
		return err
	}

	if err := balance.ReprocessFinished(correlationID); err != nil {
		return err
	}

	return nil
}
