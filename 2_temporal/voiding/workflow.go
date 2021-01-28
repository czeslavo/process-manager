package voiding

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

const (
	VoidDocumentsWorkflowName      = "void-documents-workflow"
	VoidDocumentsWorkflowQueueName = "void-documents-queue-workflow"
)

func VoidDocuments(ctx workflow.Context, documentUUID string) error {
	currentState := "initial"
	if err := workflow.SetQueryHandler(ctx, "current_state", func() (string, error) {
		return currentState, nil
	}); err != nil {
		currentState = "query_handler_register_failed"
		return err
	}

	activityOpts := workflow.ActivityOptions{
		TaskQueue:              VoidDocumentsWorkflowQueueName,
		ScheduleToCloseTimeout: time.Second * 60,
		ScheduleToStartTimeout: time.Second * 60,
		StartToCloseTimeout:    time.Second * 60,
		HeartbeatTimeout:       time.Second * 10,
		WaitForCancellation:    false,
	}
	ctx = workflow.WithActivityOptions(ctx, activityOpts)

	currentState = "marking_documents_as_voided_in_statements"
	if err := workflow.ExecuteActivity(ctx, MarkDocumentAsVoidedInReports, documentUUID).Get(ctx, nil); err != nil {
		currentState = "marking_documents_as_voided_in_statements_failed"
		return err
	}

	currentState = "voiding_documents_in_billing"
	if err := workflow.ExecuteActivity(ctx, VoidDocumentInBilling, documentUUID).Get(ctx, nil); err != nil {
		currentState = "voiding_documents_in_billing_failed"
		return err
	}

	currentState = "done"
	return nil
}
