package balance

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

func pubsub(handler *ReprocessTripHandler) (*cqrs.Facade, *message.Router) {
	logger := watermill.NewStdLogger(true, true)
	pubsub := gochannel.NewGoChannel(gochannel.Config{}, logger)

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	cqrsFacade, err := cqrs.NewFacade(cqrs.FacadeConfig{
		GenerateCommandsTopic: func(commandName string) string { return commandName },
		CommandHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.CommandHandler {
			var commandHandlers []cqrs.CommandHandler
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

	return cqrsFacade, router
}
