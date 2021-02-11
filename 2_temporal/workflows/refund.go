package workflows

import (
	"github.com/czeslavo/process-manager/2_temporal/activities"
	"github.com/czeslavo/process-manager/2_temporal/events"
	"go.temporal.io/sdk/workflow"
)

const (
	RefundWorkflowName = "refund-workflow"

	CreditNoteIssuedSignalName = "credit-note-issued-signal"
)

func Refund(ctx workflow.Context, correlationID string) error {
	logger := workflow.GetLogger(ctx)

	ctx = workflow.WithActivityOptions(ctx, activities.DefaultActivityOptions)
	future := workflow.ExecuteActivity(ctx, activities.IssueCreditNoteActivity, correlationID)
	err := future.Get(ctx, nil)
	if err != nil {
		return err
	}

	creditNoteIssued := workflow.GetSignalChannel(ctx, CreditNoteIssuedSignalName)
	var creditNoteIssuedEvent events.CreditNoteIssued
	more := creditNoteIssued.Receive(ctx, &creditNoteIssuedEvent)

	logger.Info(
		"received credit note issued event",
		"workflow_id",
		workflow.GetInfo(ctx).WorkflowExecution.ID,
		"event",
		creditNoteIssuedEvent,
		"more",
		more,
	)
	return nil
}
