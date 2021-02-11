package workflows

import "go.temporal.io/sdk/workflow"

const (
	ChargeAdditionalAmountWorkflowName = "charge-additional-amount-workflow"
)

func ChargeAdditionalAmount(ctx workflow.Context, tripUUID string) error {
	return nil
}
