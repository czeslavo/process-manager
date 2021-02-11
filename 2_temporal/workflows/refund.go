package workflows

import (
	"github.com/czeslavo/process-manager/2_temporal/events"
	"go.temporal.io/sdk/workflow"
)

const (
	RefundWorkflowName = "refund-workflow"

	CreditNoteIssuedSignalName = "credit-note-issued-signal"
)

func Refund(ctx workflow.Context, tripUUID string) error {
	creditNoteIssued := workflow.GetSignalChannel(ctx, CreditNoteIssuedSignalName)

	var creditNoteIssuedEvent events.CreditNoteIssued
	creditNoteIssued.Receive(ctx, &creditNoteIssuedEvent)

	logger := workflow.GetLogger(ctx)
	logger.Info("received credit note issued event in workflow id %s: %s", workflow.GetInfo(ctx).WorkflowExecution.ID, creditNoteIssuedEvent)

	return nil
}
