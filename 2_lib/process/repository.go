package process

import (
	"context"
	"errors"

	processmanager "github.com/czeslavo/process-manager/2_lib"
)

type Repository struct {
	storage map[string]processmanager.ProcessInstance
}

func NewRepository() *Repository {
	return &Repository{storage: make(map[string]processmanager.ProcessInstance)}
}

func (r Repository) Get(ctx context.Context, id string) (processmanager.ProcessInstance, error) {
	p, ok := r.storage[id]
	if !ok {
		return nil, processmanager.ProcessInstanceNotFound
	}

	return p, nil
}

func (r *Repository) Create(ctx context.Context, id string) (processmanager.ProcessInstance, error) {
	if _, ok := r.storage[id]; ok {
		return nil, errors.New("process already exists")
	}

	r.storage[id] = NewDocumentVoidingProcess(id)

	return r.storage[id], nil
}

func (r *Repository) Save(ctx context.Context, process processmanager.ProcessInstance) error {
	r.storage[process.ID()] = process
	return nil
}
