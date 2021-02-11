package balance

import (
	"context"

	errors "github.com/pkg/errors"
)

var repo = NewTripBalanceRepository()

func StartProcess(correlationID string) {
	// here we would call temporal workflow

}

func TriggerReprocessTrip(ctx context.Context, tripUUID, amendmentID string) error {
	balance, err := repo.GetBalance(tripUUID)
	if err != nil {
		return err
	}

	reprocessType, err := balance.ReprocessTrip(amendmentID)
	if err != nil {
		return err
	}

	switch reprocessType {
	case ReprocessTypeRefund:
		StartProcess(amendmentID)
	case ReprocessTypeAdditionalPayment:
		StartProcess(amendmentID)
	case ReprocessTypeChangedContractDetails:
		StartProcess(amendmentID)
	default:
		return errors.New("unknown reprocess type")
	}

	return nil
}
