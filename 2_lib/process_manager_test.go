package processmanager_test

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"testing"

	processmanager "github.com/czeslavo/process-manager/2_lib"
	sampleProcess "github.com/czeslavo/process-manager/2_lib/process"
	"github.com/stretchr/testify/require"
)

func TestProcessManager(t *testing.T) {
	uuidProvider := func() string {
		return "uuid"
	}

	processManager, err := processmanager.New(processmanager.Config{
		UUIDProvider:  uuidProvider,
		EventsMapping: eventsMapping,
		Repository:    sampleProcess.NewRepository(),
	})
	require.NoError(t, err)

	eventHandlers, err := processManager.EventHandlers()
	require.NoError(t, err)
	require.NotEmpty(t, eventHandlers)

	eventsSequence := []interface{}{

	}

}

type voidingRequested struct {
	DocumentID string
}

var eventsMapping = map[interface{}]func(e interface{}) (processmanager.Event, error){
	&voidingRequested{}: func(e interface{}) (processmanager.Event, error) {
		transportEvent := e.(*voidingRequested)

		return processmanager.Event{
			CorrelationID: transportEvent.DocumentID,
			Payload:
		}, nil
	},
}
