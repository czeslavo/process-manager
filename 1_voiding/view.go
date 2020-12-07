package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/czeslavo/process-manager/1_voiding/billing"
	"github.com/czeslavo/process-manager/1_voiding/manager"
	"github.com/czeslavo/process-manager/1_voiding/messages"
	"github.com/czeslavo/process-manager/1_voiding/reports"
)

func runHttpServer(
	ctx context.Context,
	commandBus *cqrs.CommandBus,
	reportsService *reports.Service,
	billingService *billing.Service,
	processManager *manager.DocumentVoidingProcessManager,
) {
	http.HandleFunc("/documents/void", func(w http.ResponseWriter, r *http.Request) {
		documentID := r.PostFormValue("id")
		if err := commandBus.Send(ctx, &messages.RequestDocumentVoiding{
			DocumentID: documentID,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	})

	http.HandleFunc("/reports/publish", func(w http.ResponseWriter, r *http.Request) {
		customerID := r.PostFormValue("id")
		if err := commandBus.Send(ctx, &messages.PublishReport{CustomerID: customerID}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	})

	http.HandleFunc("/processes/ack", func(w http.ResponseWriter, r *http.Request) {
		processID := r.PostFormValue("id")
		if err := commandBus.Send(ctx, &messages.AcknowledgeProcessFailure{ProcessID: processID}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, buildReadModel(reportsService, billingService, processManager)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func buildReadModel(reportsService *reports.Service, billingService *billing.Service, processManager *manager.DocumentVoidingProcessManager) interface{} {
	type document struct {
		ID                  string
		Amount              string
		IsVoided            bool
		VoidingProcessState string
	}

	type report struct {
		CustomerID  string
		TotalAmount string
		Documents   []document
		IsPublished bool
	}

	type process struct {
		State      string
		ID         string
		DocumentID string
	}

	type readModel struct {
		Reports   []report
		Processes []process
	}

	var reportsToDisplay []report
	for _, r := range reportsService.GetReports() {
		var docs []document
		for _, docID := range r.DocumentIDs {
			doc, err := billingService.GetDocument(docID)
			if err != nil {
				panic(err)
			}

			var processState string
			process, ongoing := processManager.GetProcessForDocument(doc.ID)
			if ongoing {
				processState = process.State.String()
			}

			docs = append(docs, document{
				ID:                  doc.ID,
				Amount:              fmt.Sprintf("%.2f", doc.TotalAmount),
				IsVoided:            doc.IsVoided,
				VoidingProcessState: processState,
			})
		}

		reportsToDisplay = append(
			reportsToDisplay, report{
				CustomerID:  r.CustomerID,
				TotalAmount: fmt.Sprintf("%.2f", r.TotalAmount),
				Documents:   docs,
				IsPublished: r.IsPublished,
			})
	}

	var processesToDisplay []process
	for _, p := range processManager.GetAllOngoingOrFailed() {
		processesToDisplay = append(
			processesToDisplay, process{
				State:      p.State.String(),
				ID:         p.ID,
				DocumentID: p.DocumentID,
			})
	}

	sort.Slice(reportsToDisplay, func(i, j int) bool { return reportsToDisplay[i].CustomerID > reportsToDisplay[j].CustomerID })
	sort.Slice(processesToDisplay, func(i, j int) bool { return processesToDisplay[i].ID > processesToDisplay[j].ID })
	return readModel{
		Reports:   reportsToDisplay,
		Processes: processesToDisplay,
	}
}
