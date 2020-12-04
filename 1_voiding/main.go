package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/czeslavo/process-manager/1_voiding/billing"
	"github.com/czeslavo/process-manager/1_voiding/manager"
	"github.com/czeslavo/process-manager/1_voiding/messages"
	"github.com/czeslavo/process-manager/1_voiding/reports"
	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()
	logger := watermill.NewStdLogger(false, false)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	billingService := billing.NewService()
	reportsService := reports.NewService()

	documentVoidingProcessManager := manager.NewDocumentVoidingProcessManager()

	cqrsFacade, err := cqrs.NewFacade(cqrs.FacadeConfig{
		GenerateCommandsTopic: func(commandName string) string { return commandName },
		CommandHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.CommandHandler {
			var commandHandlers []cqrs.CommandHandler
			commandHandlers = append(commandHandlers, billingService.CommandHandlers(eb)...)
			commandHandlers = append(commandHandlers, reportsService.CommandHandlers(eb)...)
			return commandHandlers
		},
		CommandsPublisher: pubsub,
		CommandsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			return pubsub, nil
		},
		GenerateEventsTopic: func(eventName string) string {
			return eventName
		},
		EventHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.EventHandler {
			var eventHandlers []cqrs.EventHandler
			eventHandlers = append(eventHandlers, billingService.EventHandlers()...)
			eventHandlers = append(eventHandlers, reportsService.EventHandlers()...)
			eventHandlers = append(eventHandlers, documentVoidingProcessManager.EventHandlers(cb)...)
			return eventHandlers
		},
		EventsPublisher: pubsub,
		EventsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			return pubsub, nil
		},
		Router:                router,
		CommandEventMarshaler: cqrs.JSONMarshaler{},
		Logger:                logger,
	})
	if err != nil {
		panic(err)
	}

	poisonQueueMiddleware, err := middleware.PoisonQueue(pubsub, "/dev/null")
	if err != nil {
		panic(err)
	}
	router.AddMiddleware(poisonQueueMiddleware)

	go simulateCommands(ctx, cqrsFacade.CommandBus())
	go runHttpServer(ctx, cqrsFacade.CommandBus(), reportsService, billingService, documentVoidingProcessManager)

	if err := router.Run(ctx); err != nil {
		panic(err)
	}
}

func simulateCommands(ctx context.Context, commandBus *cqrs.CommandBus) {
	recipients := []string{
		uuid.New().String(),
		uuid.New().String(),
		uuid.New().String(),
	}
	for range time.Tick(time.Second * 5) {
		if err := commandBus.Send(ctx, &messages.IssueDocument{
			DocumentID:  uuid.New().String(),
			RecipientID: recipients[rand.Intn(len(recipients))],
			TotalAmount: gofakeit.Price(1, 50),
		}); err != nil {
			fmt.Println("Failed to send command")
		}
	}
}

func runHttpServer(
	ctx context.Context,
	commandBus *cqrs.CommandBus,
	reportsService *reports.Service,
	billingService *billing.Service,
	processManager *manager.DocumentVoidingProcessManager,
) {
	http.HandleFunc("/void", func(w http.ResponseWriter, r *http.Request) {
		documentID := r.PostFormValue("id")
		if err := commandBus.Send(ctx, &messages.RequestDocumentVoiding{
			DocumentID: documentID,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		customerID := r.PostFormValue("id")
		if err := commandBus.Send(ctx, &messages.PublishReport{CustomerID: customerID}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

	type readModel struct {
		Reports []report
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

	sort.Slice(reportsToDisplay, func(i, j int) bool { return reportsToDisplay[i].CustomerID > reportsToDisplay[j].CustomerID })
	return readModel{reportsToDisplay}
}
