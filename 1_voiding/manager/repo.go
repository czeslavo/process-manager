package manager

import (
	"errors"
	"sync"
)

type Repo struct {
	sync.RWMutex
	processes map[string]DocumentVoidingProcess
}

func NewRepo() *Repo {
	return &Repo{
		processes: make(map[string]DocumentVoidingProcess),
	}
}

func (r *Repo) GetProcess(processID string) (DocumentVoidingProcess, error) {
	r.RLock()
	defer r.RUnlock()

	process, ok := r.processes[processID]
	if !ok {
		return DocumentVoidingProcess{}, errors.New("process not found")
	}
	return process, nil
}

func (r *Repo) GetOngoingForDocument(documentID string) (DocumentVoidingProcess, bool) {
	for _, p := range r.processes {
		if p.DocumentID == documentID && p.IsOngoing() {
			return p, true
		}
	}
	return DocumentVoidingProcess{}, false
}

func (r *Repo) GetAllOngoing() []DocumentVoidingProcess {
	var processes []DocumentVoidingProcess
	for _, p := range r.processes {
		if p.IsOngoing() {
			processes = append(processes, p)
		}
	}

	return processes
}

func (r *Repo) Store(process DocumentVoidingProcess) {
	r.Lock()
	defer r.Unlock()

	r.processes[process.ID] = process
}

func (r *Repo) IsOngoingForDocument(documentID string) bool {
	r.RLock()
	defer r.RUnlock()

	for _, p := range r.processes {
		if p.DocumentID == documentID && p.IsOngoing() {
			return true
		}
	}
	return false
}
