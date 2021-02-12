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

func (c temporal) startWorkflow(
	ctx context.Context,
	workflowParams WorkflowParams,
	reprocessType ReprocessType,
) error {
	workflowName, err := getWorkflowName(reprocessType)

	run, err := c.client.ExecuteWorkflow(
		ctx,
		client.StartWorkflowOptions{
			ID:        workflowParams.CorrelationID,
			TaskQueue: worker.TaskQueue,
		},
		workflowName,
		workflowParams,
	)
	if err != nil {
		return err
	}

	c.logger.
		WithField("correlation_id", workflowParams.CorrelationID).
		WithField("run_id", run.GetRunID()).
		WithField("reprocess_type", reprocessType).
		Info("Run started")
	return nil
}
