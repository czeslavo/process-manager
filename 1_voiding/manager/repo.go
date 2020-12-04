package manager

import "sync"

type Repo struct {
	sync.RWMutex
	processes map[string]DocumentVoidingProcess
}

func NewRepo() *Repo {
	return &Repo{
		processes: make(map[string]DocumentVoidingProcess),
	}
}

func (r *Repo) GetOrCreateProcess(processID, documentID string) DocumentVoidingProcess {
	r.Lock()
	defer r.Unlock()

	if process, ok := r.processes[processID]; !ok || !process.IsOngoing() {
		return NewDocumentVoidingProcess(processID, documentID)
	}

	return r.processes[processID]
}

func (r *Repo) GetOngoingForDocument(documentID string) (DocumentVoidingProcess, bool) {
	for _, p := range r.processes {
		if p.DocumentID == documentID && p.IsOngoing() {
			return p, true
		}
	}
	return DocumentVoidingProcess{}, false
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
