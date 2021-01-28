package reports

import "sync"

type Report struct {
	CustomerID  string
	TotalAmount float64
	DocumentIDs []string
	IsPublished bool
}

func (r *Report) AppendDocument(documentID string, documentTotal float64) {
	r.TotalAmount = r.TotalAmount + documentTotal
	r.DocumentIDs = append(r.DocumentIDs, documentID)
}

func (r *Report) Publish() {
	r.IsPublished = true
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

func (r *Repo) GetOrCreate(customerID string) Report {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.customerReport[customerID]; !ok {
		return Report{
			CustomerID: customerID,
		}
	}

	return r.customerReport[customerID]
}

func (r *Repo) GetForDocument(documentID string) []Report {
	r.Lock()
	defer r.Unlock()

	var reports []Report

	for _, report := range r.customerReport {
		for _, inReportDocumentID := range report.DocumentIDs {
			if inReportDocumentID == documentID {
				reports = append(reports, report)
			}
		}
	}

	return reports
}

func (r *Repo) Store(report Report) {
	r.Lock()
	defer r.Unlock()

	r.customerReport[report.CustomerID] = report
}
