package stream

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

type Message struct {
	typ      Name
	sequence Sequence

	// body TODO
	body Event

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

func NewMessage(s Sequence, e Event) Message {
	return Message{
		sequence:  s,
		typ:       Name(reflect.TypeOf(e).Name()),
		body:      e,
		createdAt: time.Now(),
	}
}

func (m Message) ID() string {
	return m.sequence.ID()
}

func (m Message) Sequence() Sequence {
	return m.sequence
}

func (m Message) Correlate(d ID) Message {
	m.correlation = d
	return m
}

func (m Message) Respond(src Message) Message {
	m.correlation, m.causation = src.correlation, src.sequence.namespace.id
	return m
}

func (m Message) String() string {
	return fmt.Sprintf("%s.%s#%d", m.sequence.namespace.id, m.typ, m.sequence.number)
}

func (m Message) GoString() string {
	v := view{
		"ID":          m.sequence.ID(),
		"Type":        m.sequence.namespace.name + m.typ,
		"Correlation": m.correlation,
		"Causation":   m.causation,
		"Namespace":   m.sequence.namespace,
		"CreatedAt":   m.createdAt,
		"Body":        m.body,
		"Meta":        m.meta,
	}
	b, _ := json.MarshalIndent(v, "", "\t")
	return fmt.Sprintf("%T\n%s\n", m, b)
}
