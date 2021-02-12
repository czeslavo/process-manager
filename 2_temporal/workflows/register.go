package workflows

import (
	"github.com/czeslavo/process-manager/2_temporal/balance"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func RegisterWorkflows(w worker.Worker) error {
	w.RegisterWorkflow(ChargeAdditionalAmount)
	w.RegisterWorkflowWithOptions(Refund, workflow.RegisterOptions{
		Name:                          balance.RefundWorkflowName,
		DisableAlreadyRegisteredCheck: false,
	})
	w.RegisterWorkflow(ReissueDocuments)

	return nil
}
