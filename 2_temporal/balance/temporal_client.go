package balance

import (
	"context"

	"github.com/czeslavo/process-manager/2_temporal/worker"

	"github.com/sirupsen/logrus"
	"go.temporal.io/sdk/client"
)

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
			TaskQueue: worker.TaskQueue,
		},
		workflowName,
		correlationID,
	)
	if err != nil {
		return err
	}

	c.logger.
		WithField("correlation_id", correlationID).
		WithField("run_id", run.GetRunID()).
		WithField("reprocess_type", reprocessType).
		Info("Run started")
	return nil
}
