package balance

import errors "github.com/pkg/errors"

type ReprocessType string

var (
	ReprocessTypeRefund                 = ReprocessType("refund")
	ReprocessTypeAdditionalPayment      = ReprocessType("additional-payment")
	ReprocessTypeChangedContractDetails = ReprocessType("changed-contract-details")
)

type TripBalance struct {
	amountToBePaid float64
	tripUUID       string

	amendmentInProgress string
}

func NewTripBalance(tripUUID string, b float64) TripBalance {
	return TripBalance{b, tripUUID, ""}
}

func (b TripBalance) ReprocessTrip(amendmentID string) (ReprocessType, error) {
	if b.amendmentInProgress != "" {
		return "", errors.Errorf("cannot start new reprocess, amendment in progress: %s", b.amendmentInProgress)
	}

	b.amendmentInProgress = amendmentID

	if b.amountToBePaid > 0 {
		return ReprocessTypeAdditionalPayment, nil
	}
	if b.amountToBePaid < 0 {
		return ReprocessTypeRefund, nil
	}
	if b.amountToBePaid == 0 {
		return ReprocessTypeChangedContractDetails, nil
	}

	return "", errors.New("dunno what to do")
}

func (b *TripBalance) ReprocessFinished(amendmentID string) error {
	if b.amendmentInProgress != amendmentID {
		return errors.New("cannot finish process because amendment ids don't match")
	}

	b.amendmentInProgress = ""
	return nil
}

func (b TripBalance) TripUUID() string {
	return b.tripUUID
}
