package workflows

import "go.temporal.io/sdk/workflow"

const ReissueDocumentsWorkflowName = "reissue-documents-workflow"

func ReissueDocuments(ctx workflow.Context, tripUUID string) error {
	return nil
}
