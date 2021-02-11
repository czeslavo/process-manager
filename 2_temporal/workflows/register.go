package workflows

import (
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func RegisterWorkflows(w worker.Worker) error {
	w.RegisterWorkflow(ChargeAdditionalAmount)
	w.RegisterWorkflowWithOptions(Refund, workflow.RegisterOptions{
		Name:                          RefundWorkflowName,
		DisableAlreadyRegisteredCheck: false,
	})
	w.RegisterWorkflow(ReissueDocuments)

	return nil
}
