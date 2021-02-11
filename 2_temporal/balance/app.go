package balance

import (
	"context"
	"encoding/json"

	"github.com/czeslavo/process-manager/2_temporal/events"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/czeslavo/process-manager/2_temporal/workflows"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.temporal.io/sdk/client"
)

type ReprocessTripHandler struct {
	repo           *TripBalanceRepository
	temporalClient *temporal
	logger         logrus.FieldLogger
	publisher      message.Publisher
}

func NewReprocessTripHandler(repo *TripBalanceRepository, publisher message.Publisher) (*ReprocessTripHandler, error) {
	logger := logrus.New()
	c, err := client.NewClient(client.Options{})
	if err != nil {
		return nil, err
	}

	t := &temporal{
		client: c,
		logger: logger,
	}
	return &ReprocessTripHandler{
		repo:           repo,
		logger:         logger,
		temporalClient: t,
		publisher:      publisher,
	}, nil
}

func (h ReprocessTripHandler) HandleReprocess(ctx context.Context, tripUUID, correlationID string) error {
	balance, err := h.repo.GetBalance(tripUUID)
	if err != nil {
		return err
	}

	reprocessType, err := balance.ReprocessTrip(correlationID)
	if err != nil {
		return err
	}

	return h.temporalClient.startWorkflow(ctx, correlationID, reprocessType)
}

func (h ReprocessTripHandler) HandleReprocessFinished(ctx context.Context, tripUUID, correlationID string) error {
	balance, err := h.repo.GetBalance(tripUUID)
	if err != nil {
		return err
	}

	if err := balance.ReprocessFinished(correlationID); err != nil {
		return err
	}

	return nil
}

func (h ReprocessTripHandler) HandleCreditNoteIssued(msg *message.Message) error {
	var event events.CreditNoteIssued
	if err := json.Unmarshal(msg.Payload, &event); err != nil {
		return err
	}

	if err := h.temporalClient.client.SignalWorkflow(
		msg.Context(),
		event.CorrelationID,
		"", // will use latest running workflow for this correlation id (in our case it's always one or none)
		workflows.CreditNoteIssuedSignalName,
		event,
	); err != nil {
		return err
	}

	return nil
}

func getWorkflowName(reprocessType ReprocessType) (string, error) {
	switch reprocessType {
	case ReprocessTypeRefund:
		return workflows.RefundWorkflowName, nil
	case ReprocessTypeAdditionalPayment:
		return workflows.ChargeAdditionalAmountWorkflowName, nil
	case ReprocessTypeChangedContractDetails:
		return workflows.ReissueDocumentsWorkflowName, nil
	default:
		return "", errors.New("unsupported reprocess type")
	}
}
