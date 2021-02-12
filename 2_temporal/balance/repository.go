package balance

import (
	"sync"

	errors "github.com/pkg/errors"
)

type TripBalanceRepository struct {
	mtx      *sync.RWMutex
	balances map[string]TripBalance
}

func NewTripBalanceRepository() *TripBalanceRepository {
	return &TripBalanceRepository{balances: make(map[string]TripBalance), mtx: &sync.RWMutex{}}
}

func (r TripBalanceRepository) GetBalance(tripUUID string) (TripBalance, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	b, ok := r.balances[tripUUID]
	if !ok {
		return TripBalance{}, errors.New("missing trip balance")
	}

	return b, nil
}

func (r *TripBalanceRepository) SaveBalance(balance TripBalance) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.balances[balance.TripUUID()] = balance
}
