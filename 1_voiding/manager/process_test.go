package manager_test

import (
	"testing"

	"github.com/czeslavo/process-manager/1_voiding/manager"
	"github.com/czeslavo/process-manager/1_voiding/messages"
	"github.com/stretchr/testify/require"
)

func TestDocumentVoidingProcess_correct_scenarios(t *testing.T) {
	t.Run("documents_voided", func(t *testing.T) {
		p := spawnProcess()
		nextCommand := p.NextCommand()
		require.IsType(t, &messages.MarkDocumentAsVoided{}, nextCommand)

		err := p.MarkingDocumentAsVoidedSucceeded()
		require.NoError(t, err)

		nextCommand = p.NextCommand()
		require.IsType(t, &messages.VoidDocument{}, nextCommand)

		err = p.DocumentVoided()
		require.NoError(t, err)

		nextCommand = p.NextCommand()
		require.Nil(t, nextCommand)
	})

	t.Run("marking_documents_as_voided_failed", func(t *testing.T) {
		p := spawnProcess()
		nextCommand := p.NextCommand()
		require.IsType(t, &messages.MarkDocumentAsVoided{}, nextCommand)

		err := p.MarkingDocumentAsVoidedFailed()
		require.NoError(t, err)

		nextCommand = p.NextCommand()
		require.Nil(t, nextCommand)

		err = p.AcknowledgeFailure()
		require.NoError(t, err)

	})
}

func spawnProcess() manager.DocumentVoidingProcess {
	processID := "process-id"
	documentID := "document-id"
	customerID := "customer-id"
	return manager.NewDocumentVoidingProcess(processID, documentID, customerID)
}
