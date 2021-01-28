package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"go.temporal.io/sdk/client"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/brianvoe/gofakeit/v5"
	"github.com/czeslavo/process-manager/2_temporal/billing"
	"github.com/czeslavo/process-manager/2_temporal/messages"
	"github.com/czeslavo/process-manager/2_temporal/reports"
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

	temporalClient, err := client.NewClient(client.Options{})
	if err != nil {
		panic(err)
	}
	defer temporalClient.Close()

	billingService := billing.NewService(temporalClient)
	reportsService := reports.NewService()

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
	router.AddMiddleware(slowDown)

	go simulateCommands(ctx, cqrsFacade.CommandBus())
	go runHttpServer(ctx, cqrsFacade.CommandBus(), reportsService, billingService)

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

func slowDown(h message.HandlerFunc) message.HandlerFunc {
	shouldSlowDown := os.Getenv("SLOW_DOWN") == "1"
	return func(message *message.Message) ([]*message.Message, error) {
		if shouldSlowDown {
			time.Sleep(time.Millisecond * 250)
		}

		return h(message)
	}
}
