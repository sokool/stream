package stream

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Event[E any] struct {
	typ       Type
	namespace Namespace
	sequence  int64
	// body TODO
	body E

	// meta TODO
	meta Meta

	// Every Message has 3 ID's [ID, CorrelationID, CausationID]. When you are
	// responding to a Message (either a Command or and Event) you copy the
	// CorrelationID of the Message you are responding to, to your new
	// CorrelationID. The CausationID of your Message is the ID of the
	// Message you are responding to.
	//
	// Greg Young
	// --> https://groups.google.com/d/msg/dddcqrs/qGYC6qZEqOI/LhQup9v7EwAJ
	correlation, causation ID

	// CreatedAt
	createdAt time.Time

	// author helps to check what person/device generate this Message.
	author string
}

func NewEvent[E any](n Namespace, e E, sequence int64) (m Event[E], err error) {
	m = Event[E]{
		namespace:   n,
		sequence:    sequence,
		body:        e,
		meta:        Meta{},
		correlation: "",
		causation:   "",
		createdAt:   time.Now(),
		author:      "",
	}

	if m.typ, err = NewType(e); err != nil {
		return m, nil
	}

	return m, nil
}

//func DecodeEvent[E any]()

func (e Event[E]) ID() ID {
	return ID(uuid.NewSHA1(uuid.NameSpaceDNS, []byte(e.String())).String())
}

func (e Event[E]) Type() Type {
	return Type("")
}

func (e Event[E]) Namespace() Namespace {
	return e.namespace
}

func (e Event[E]) Sequence() int64 {
	return e.sequence
}

func (e Event[E]) Body() E {
	return e.body
}

func (e Event[E]) Correlate(d ID) Event[E] {
	e.correlation = d
	return e
}

func (e Event[E]) Respond(src Event[E]) Event[E] {
	e.correlation, e.causation = src.correlation, src.ID()
	return e
}

func (e Event[E]) String() string {
	return fmt.Sprintf("%s.%s#%d", e.namespace, e.typ, e.sequence)
}

func (e Event[E]) GoString() string {
	v := view{
		"ID":          e.ID(),
		"Type":        e.namespace.name + e.typ,
		"Correlation": e.correlation,
		"Causation":   e.causation,
		"Namespace":   e.namespace,
		"CreatedAt":   e.createdAt,
		"Body":        e.body,
		"Meta":        e.meta,
	}
	b, _ := json.MarshalIndent(v, "", "\t")
	return fmt.Sprintf("%T\n%s\n", e, b)
}
