package activities

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/czeslavo/process-manager/2_temporal/balance"
)

type Handler struct {
	publisher message.Publisher
	repo      *balance.TripBalanceRepository
}

func NewHandler(publisher message.Publisher, repo *balance.TripBalanceRepository) Handler {
	return Handler{publisher: publisher, repo: repo}
}
