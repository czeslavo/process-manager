package balance

import (
	"context"

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
