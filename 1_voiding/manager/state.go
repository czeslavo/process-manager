package manager

import "github.com/pkg/errors"

/*
State machine describing the process:

 +------------------+      +-------------------------------+          +----------------------+
 | VoidingRequested | +--> | MarkingDocumentAsVoidedFailed | +------> | Failure Acknowledged |
 +------------------+      +-------------------------------+          +----------------------+
         |
         +                 +----------------------------------+          +----------------+
         +---------------> | MarkingDocumentAsVoidedSucceeded | +------> | DocumentVoided |
                           +----------------------------------+          +----------------+
*/

type State string

const (
	VoidingRequested                 State = "voiding-requested"
	MarkingDocumentAsVoidedFailed    State = "marking-document-as-voided-failed"
	MarkingDocumentAsVoidedSucceeded State = "marking-document-as-voided-succeeded"
	DocumentVoided                   State = "document-voided"
	FailureAcknowledged              State = "failure-acknowledged"
)

func (s State) String() string {
	return string(s)
}

func (s State) canTransition(to State) error {
	allowedStateTransitions := allowedTransitions{
		VoidingRequested:                 []State{MarkingDocumentAsVoidedSucceeded, MarkingDocumentAsVoidedFailed},
		MarkingDocumentAsVoidedSucceeded: []State{DocumentVoided},
		MarkingDocumentAsVoidedFailed:    []State{FailureAcknowledged},
		DocumentVoided:                   nil,
		FailureAcknowledged:              nil,
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
