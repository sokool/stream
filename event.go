package stream

import (
	"encoding/json"
	"fmt"
	"time"
)

type Event[T any] struct {
	ID       ID
	Type     Type
	Root     RootID
	Sequence int64

	// body TODO
	Body T

	// meta TODO
	Meta Meta

	// Every Message has 3 ID's [ID, CorrelationID, CausationID]. When you are
	// responding to a Message (either a Command or and Event) you copy the
	// CorrelationID of the Message you are responding to, to your new
	// CorrelationID. The CausationID of your Message is the ID of the
	// Message you are responding to.
	//
	// Greg Young
	// --> https://groups.google.com/d/msg/dddcqrs/qGYC6qZEqOI/LhQup9v7EwAJ
	Correlation, Causation ID

	// CreatedAt
	CreatedAt time.Time

	// Author helps to check what person/device generate this Message.
	Author string
}

func NewEvent[T any](id RootID, v T, sequence int64) (e Event[T], err error) {
	var t Type
	if t, err = NewType(v); err != nil {
		return e, nil
	}

	if sequence <= 0 {
		return e, Err("invalid event sequence")
	}

	return Event[T]{
		ID:          uid(fmt.Sprintf("%s.%s.%d", id, e.Type, sequence)),
		Root:        id,
		Type:        t.CutPrefix(id.Type()),
		Sequence:    sequence,
		Body:        v,
		Meta:        Meta{},
		Correlation: "",
		Causation:   "",
		CreatedAt:   time.Now(),
		Author:      "",
	}, nil
}

func (e *Event[T]) Correlate(d ID) *Event[T] {
	e.Correlation = d
	return e
}

func (e *Event[T]) Respond(to Event[any]) *Event[T] {
	e.Correlation, e.Causation = to.Correlation, to.ID
	return e
}

func (e *Event[E]) String() string {
	return fmt.Sprintf("%s%s#%d", e.Root, e.Type, e.Sequence)
}

func (e *Event[T]) GoString() string {
	b, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("%T\n%s\n", e.Body, b)
}

func (e *Event[T]) IsZero() bool {
	return e.ID == ""
}

type Events []Event[any]

func NewEvents(r Root) (ee Events, err error) {
	var id RootID
	var e Event[any]
	var k = r.Version() + 1
	if id, err = NewRootID(r); err != nil {
		return nil, err
	}

	for i, v := range r.Uncommitted(true) {
		if e, err = NewEvent(id, v, k+int64(i)); err != nil {
			return nil, err
		}
		ee = append(ee, e)
	}

	return ee, nil
}

// Unique gives RootID when all events has same RootID
func (e Events) Unique() RootID {
	if len(e) == 0 || !e.hasUnique(e[0].Root) {
		return RootID{}
	}

	return e[0].Root
}

func (e Events) hasUnique(id RootID) bool {
	for i := range e {
		if id != e[i].Root && !e[i].IsZero() {
			return false
		}
	}
	return true
}

func (e Events) String() string {
	var s string
	for i := range e {
		if e[i].IsZero() {
			continue
		}
		s += fmt.Sprintf("%s\n", &e[i])
	}

	return s
}

//func DecodeEvent[E any]()
