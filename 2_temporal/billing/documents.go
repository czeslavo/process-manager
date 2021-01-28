package billing

import (
	"errors"
	"sync"
)

type Document struct {
	ID          string
	RecipientID string
	TotalAmount float64 // float64 for simplicity
	IsVoided    bool
}

func (d *Document) Void() {
	d.IsVoided = true
}

type DocumentsRepo struct {
	sync.RWMutex
	documents map[string]Document
}

func NewDocumentsRepo() *DocumentsRepo {
	return &DocumentsRepo{
		documents: make(map[string]Document),
	}
}

func (r *DocumentsRepo) GetByID(id string) (Document, error) {
	r.RLock()
	defer r.RUnlock()

	d, ok := r.documents[id]
	if !ok {
		return Document{}, errors.New("document not found")
	}

	return d, nil
}

func (r *DocumentsRepo) Store(id string, doc Document) {
	r.Lock()
	defer r.Unlock()

	r.documents[id] = doc
}
