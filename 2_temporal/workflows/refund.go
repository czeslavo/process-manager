package workflows

import "go.temporal.io/sdk/workflow"

const (
	RefundWorkflowName = "refund-workflow"
)

func Refund(ctx workflow.Context, tripUUID string) error {

	return nil
}
