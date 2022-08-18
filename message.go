package stream

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Message[E any] struct {
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

func NewMessage[E any](n Namespace, e E, sequence int64) (m Message[E], err error) {
	m = Message[E]{
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

func (m Message[E]) ID() ID {
	return ID(uuid.NewSHA1(uuid.NameSpaceDNS, []byte(m.String())).String())
}

func (m Message[E]) Sequence() int64 {
	return m.sequence
}

func (m Message[E]) Correlate(d ID) Message[E] {
	m.correlation = d
	return m
}

func (m Message[E]) Respond(src Message[E]) Message[E] {
	m.correlation, m.causation = src.correlation, src.ID()
	return m
}

func (m Message[E]) String() string {
	return fmt.Sprintf("%s.%s#%d", m.namespace, m.typ, m.sequence)
}

func (m Message[E]) GoString() string {
	v := view{
		"ID":          m.ID(),
		"Type":        m.namespace.name + m.typ,
		"Correlation": m.correlation,
		"Causation":   m.causation,
		"Namespace":   m.namespace,
		"CreatedAt":   m.createdAt,
		"Body":        m.body,
		"Meta":        m.meta,
	}
	b, _ := json.MarshalIndent(v, "", "\t")
	return fmt.Sprintf("%T\n%s\n", m, b)
}
