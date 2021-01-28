package voiding_test

import (
	"context"
	"testing"

	"github.com/czeslavo/process-manager/2_temporal/voiding"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type VoidDocumentsWorkflowTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *VoidDocumentsWorkflowTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *VoidDocumentsWorkflowTestSuite) AfterTest(_, _ string) {
	s.env.AssertExpectations(s.T())
}

func (s *VoidDocumentsWorkflowTestSuite) TestVoidDocumentsWorkflow() {
	expectedDocumentUUID := "my-document-uuid"

	s.env.OnActivity(voiding.MarkDocumentAsVoidedInReports, mock.Anything, mock.Anything).Return(func(ctx context.Context, documentUUID string) error {
		return nil
	})

	s.env.OnActivity(voiding.VoidDocumentInBilling, mock.Anything, mock.Anything).Return(func(ctx context.Context, documentUUID string) error {
		return nil
	})

	s.env.ExecuteWorkflow(voiding.VoidDocuments, expectedDocumentUUID)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}

func TestVoidDocuments(t *testing.T) {
	suite.Run(t, &VoidDocumentsWorkflowTestSuite{})
}
