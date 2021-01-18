package processmanager

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type ProcessManager struct {
	uuidProvider func() string
	config       Config
	repository   ProcessRepository
	commandBus   *cqrs.CommandBus
}

type Config struct {
	UUIDProvider  func() string
	EventsMapping map[interface{}]func(e interface{}) (Event, error)
	Repository    ProcessRepository
}

type Event struct {
	CorrelationID string
	Payload       interface{}
}

type Command interface{}

type ProcessRepository interface {
	Get(ctx context.Context, id string) (ProcessInstance, error)
	Create(ctx context.Context, id string) (ProcessInstance, error)
	Save(ctx context.Context, process ProcessInstance) error
}

type ProcessInstance interface {
	Process(event Event) error
	NextCommands() []interface{} // todo: define type
	ID() string
	State() string
}

var ProcessInstanceNotFound = errors.New("process instance not found")

func New(cfg Config) (*ProcessManager, error) {
	if cfg.UUIDProvider == nil {
		return nil, errors.New("missing uuid provider")
	}
	if len(cfg.EventsMapping) == 0 {
		return nil, errors.New("missing events mapping")
	}
	if cfg.Repository == nil {
		return nil, errors.New("missing repository")
	}

	return &ProcessManager{
		uuidProvider: cfg.UUIDProvider,
		config:       cfg,
		repository:   cfg.Repository,
	}, nil
}

func (m *ProcessManager) EventHandlers() ([]cqrs.EventHandler, error) {
	var handlers []cqrs.EventHandler

	for eventType, mapEventToDomain := range m.config.EventsMapping {
		typeOfEvent := reflect.TypeOf(eventType)

		if typeOfEvent.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("event mapping keys should be pointers: %s", typeOfEvent.String())
		}

		handlers = append(handlers, EventHandlerFunc(
			typeOfEvent.Name(),
			eventType,
			func(ctx context.Context, event interface{}) error {
				domainEvent, err := mapEventToDomain(event)
				if err != nil {
					return fmt.Errorf("could not map transport to domain event: %w", err)
				}

				p, err := m.repository.Get(ctx, domainEvent.CorrelationID)
				if errors.Is(err, ProcessInstanceNotFound) {
					// todo: only if is initial event??
					processID := m.uuidProvider()
					p, err = m.repository.Create(ctx, processID)
					if err != nil {
						return fmt.Errorf("could not create process instance")
					}
				} else if err != nil {
					return fmt.Errorf("failed to get process instance: %w", err)
				}

				if err := p.Process(domainEvent); err != nil {
					return fmt.Errorf("failed to process event: %w", err)
				}

				commands := p.NextCommands()
				_ = commands
				// map commands emitted by process instance
				// emit commands

				for _, cmd := range commands {
					// todo: what about already issues commands?
					if err := m.commandBus.Send(ctx, cmd); err != nil {
						return fmt.Errorf("failed to send command: %w", err)
					}
				}

				if err := m.repository.Save(ctx, p); err != nil {
					return fmt.Errorf("failed to save process instance: %w", err)
				}

				return nil
			},
		))
	}

	return handlers, nil
}
