package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/czeslavo/process-manager/2_temporal/billing"
	"github.com/czeslavo/process-manager/2_temporal/messages"
	"github.com/czeslavo/process-manager/2_temporal/reports"
)

func runHttpServer(
	ctx context.Context,
	commandBus *cqrs.CommandBus,
	reportsService *reports.Service,
	billingService *billing.Service,
) {
	http.HandleFunc("/documents/triggerVoiding", func(w http.ResponseWriter, r *http.Request) {
		documentID := r.PostFormValue("id")

		if err := billingService.TriggerDocumentVoidingWorkflow(r.Context(), documentID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	})

	http.HandleFunc("/documents/void", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.PostForm)

		documentID := r.PostFormValue("id")

		if err := billingService.VoidDocument(r.Context(), documentID); err != nil {
			fmt.Println("error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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

	http.HandleFunc("/reports/markDocumentAsVoided", func(w http.ResponseWriter, r *http.Request) {
		documentID := r.PostFormValue("id")

		if err := reportsService.MarkDocumentAsVoided(r.Context(), documentID); err != nil {
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

		if err := tmpl.Execute(w, buildReadModel(reportsService, billingService)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func buildReadModel(reportsService *reports.Service, billingService *billing.Service) interface{} {
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

			// todo: fetch current state from temporal
			var processState string

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

	// todo: fetch all ongoing processes from temporal
	var processesToDisplay []process

	sort.Slice(reportsToDisplay, func(i, j int) bool { return reportsToDisplay[i].CustomerID > reportsToDisplay[j].CustomerID })
	sort.Slice(processesToDisplay, func(i, j int) bool { return processesToDisplay[i].ID > processesToDisplay[j].ID })
	return readModel{
		Reports:   reportsToDisplay,
		Processes: processesToDisplay,
	}
}
