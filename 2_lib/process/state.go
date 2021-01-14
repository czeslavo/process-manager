package process

import "github.com/pkg/errors"

/*
State machine describing the process:

 +------------------+      +-------------------------------+          +----------------------+
 | VoidingRequestedState | +--> | MarkingDocumentAsVoidedFailedState | +------> | Failure Acknowledged |
 +------------------+      +-------------------------------+          +----------------------+
         |
         +                 +----------------------------------+          +----------------+
         +---------------> | MarkingDocumentAsVoidedSucceededState | +------> | DocumentVoidedState |
                           +----------------------------------+          +----------------+
*/

type State string

const (
	VoidingRequestedState                 State = "voiding-requested"
	MarkingDocumentAsVoidedFailedState    State = "marking-document-as-voided-failed"
	MarkingDocumentAsVoidedSucceededState State = "marking-document-as-voided-succeeded"
	DocumentVoidedState                   State = "document-voided"
	FailureAcknowledgedState              State = "failure-acknowledged"
)

func (s State) String() string {
	return string(s)
}

func (s State) canTransition(to State) error {
	allowedStateTransitions := allowedTransitions{
		VoidingRequestedState: []State{
			MarkingDocumentAsVoidedSucceededState,
			MarkingDocumentAsVoidedFailedState,
		},
		MarkingDocumentAsVoidedSucceededState: []State{DocumentVoidedState},
		MarkingDocumentAsVoidedFailedState:    []State{FailureAcknowledgedState},
		DocumentVoidedState:                   nil,
		FailureAcknowledgedState:              nil,
	}

	if !allowedStateTransitions.allowed(s, to) {
		return errors.Errorf("transition from '%s' to '%s' not allowed", s, to)
	}

	return nil
}

type allowedTransitions map[State][]State

func (t allowedTransitions) allowed(from, to State) bool {
	allowedTos, ok := t[from]
	if !ok {
		return false
	}

	for _, allowedTo := range allowedTos {
		if to == allowedTo {
			return true
		}
	}

	return false
}
