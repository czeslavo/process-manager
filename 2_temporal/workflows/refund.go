package workflows

import (
	"github.com/czeslavo/process-manager/2_temporal/activities"
	"github.com/czeslavo/process-manager/2_temporal/balance"
	"github.com/czeslavo/process-manager/2_temporal/events"
	"go.temporal.io/sdk/workflow"
)

func Refund(ctx workflow.Context, workflowParams balance.WorkflowParams) error {
	logger := workflow.GetLogger(ctx)

	ctx = workflow.WithActivityOptions(ctx, activities.DefaultActivityOptions)
	workflow.ExecuteActivity(ctx, activities.IssueCreditNoteActivity, workflowParams.CorrelationID)
	workflow.ExecuteActivity(ctx, activities.TriggerRefundActivity, workflowParams.CorrelationID)

	creditNoteIssued := workflow.GetSignalChannel(ctx, balance.CreditNoteIssuedSignalName)
	var creditNoteIssuedEvent events.CreditNoteIssued

	refunded := workflow.GetSignalChannel(ctx, balance.RefundedSignalName)
	var refundedEvent events.Refunded

	wg := workflow.NewWaitGroup(ctx)
	wg.Add(2)

	workflow.Go(ctx, func(ctx workflow.Context) {
		creditNoteIssued.Receive(ctx, &creditNoteIssuedEvent)
		logger.Info("Received signal for credit note issued")
		wg.Done()
	})
	workflow.Go(ctx, func(ctx workflow.Context) {
		refunded.Receive(ctx, &refundedEvent)
		logger.Info("Received signal for refunded")
		wg.Done()
	})

	if err := workflow.ExecuteActivity(ctx, activities.ReprocessingFinishedActivity, workflowParams.TripUUID, workflowParams.CorrelationID).Get(ctx, nil); err != nil {
		return err
	}

	wg.Wait(ctx)
	return nil
}
