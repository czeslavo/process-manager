package messages

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

func EventHandlerFunc(
	name string,
	eventType interface{},
	handlerFunc func(ctx context.Context, event interface{}) error,
) cqrs.EventHandler {
	return handler{
		name:        name,
		messageType: eventType,
		handlerFunc: handlerFunc,
	}
}

func CommandHandlerFunc(
	name string,
	eventType interface{},
	handlerFunc func(ctx context.Context, event interface{}) error,
) cqrs.CommandHandler {
	return handler{
		name:        name,
		messageType: eventType,
		handlerFunc: handlerFunc,
	}
}

// handler conforms cqrs.CommandHandler and cqrs.EventHandler interfaces
type handler struct {
	name        string
	messageType interface{}
	handlerFunc func(ctx context.Context, event interface{}) error
}

func (h handler) HandlerName() string {
	return h.name
}

func (h handler) NewEvent() interface{} {
	return h.messageType
}

func (h handler) NewCommand() interface{} {
	return h.messageType
}

func (h handler) Handle(ctx context.Context, event interface{}) error {
	return h.handlerFunc(ctx, event)
}
