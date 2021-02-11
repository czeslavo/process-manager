package balance

import errors "github.com/pkg/errors"

type TripBalanceRepository struct {
	balances map[string]TripBalance
}

func NewTripBalanceRepository() *TripBalanceRepository {
	return &TripBalanceRepository{balances: make(map[string]TripBalance)}
}

func (r TripBalanceRepository) GetBalance(tripUUID string) (TripBalance, error) {
	b, ok := r.balances[tripUUID]
	if !ok {
		return TripBalance{}, errors.New("missing trip balance")
	}

	return b, nil
}

func (r *TripBalanceRepository) SaveBalance(balance TripBalance) {
	r.balances[balance.TripUUID()] = balance
}
