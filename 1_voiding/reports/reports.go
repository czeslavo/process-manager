package reports

import "sync"

type Report struct {
	CustomerID  string
	TotalAmount float64
	DocumentIDs []string
}

type Repo struct {
	sync.RWMutex
	customerReport map[string]Report
}

func NewRepo() *Repo {
	return &Repo{
		customerReport: make(map[string]Report),
	}
}

func (r *Repo) AppendToReport(customerID, documentID string, documentTotal float64) {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.customerReport[customerID]; !ok {
		r.customerReport[customerID] = Report{
			CustomerID:  customerID,
			TotalAmount: documentTotal,
			DocumentIDs: []string{documentID},
		}
		return
	}

	report := r.customerReport[customerID]
	report.TotalAmount = report.TotalAmount + documentTotal
	report.DocumentIDs = append(report.DocumentIDs, documentID)

	r.customerReport[customerID] = report
}
