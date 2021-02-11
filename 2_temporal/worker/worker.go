package worker

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const TaskQueue = "trip_balance"

func NewWorker() (worker.Worker, error) {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		return nil, err
	}

	w := worker.New(
		c,
		TaskQueue,
		worker.Options{},
	)
	return w, nil
}
