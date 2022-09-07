package stream

import (
	"encoding/json"
	"fmt"
	"time"
)

type Event[E any] struct {
	typ      Type
	root     Namespace
	sequence int64

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
		root:        n,
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

	if sequence <= 0 {
		return m, Err("invalid event sequence")
	}

	return m, nil
}

func (e Event[E]) ID() ID {
	return uid(e)
}

func (e Event[E]) Type() Type {
	return e.typ
}

func (e Event[E]) Namespace() Namespace {
	return e.root
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

func (e Event[E]) Respond(src Event[any]) Event[E] {
	e.correlation, e.causation = src.correlation, src.ID()
	return e
}

func (e Event[E]) String() string {
	return fmt.Sprintf("%s.%s#%d", e.root, e.typ, e.sequence)
}

func (e Event[E]) GoString() string {
	v := view{
		"ID":          e.ID(),
		"Type":        e.root.root + e.typ,
		"Correlation": e.correlation,
		"Causation":   e.causation,
		"Namespace":   e.root,
		"CreatedAt":   e.createdAt,
		"Body":        e.body,
		"Meta":        e.meta,
	}
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("%T\n%s\n", e, b)
}

func (e Event[E]) MarshalJSON() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (e Event[E]) UnmarshalJSON(b []byte) (err error) {
	fmt.Println(string(b))
	var event struct {
		Type      Type
		Namespace Namespace
		Sequence  int64
	}

	if err = json.Unmarshal(b, &event); err != nil {
		return err
	}

	e.typ = event.Type
	e.sequence = event.Sequence
	e.root = event.Namespace

	return nil
}

type Events []Event[any]

func NewEvents(r Root) (ee Events, err error) {
	var n Namespace
	var e Event[any]
	var v = r.Version() + 1
	if n, err = NewNamespace(r); err != nil {
		return nil, err
	}

	for i, m := range r.Uncommitted(true) {
		if e, err = NewEvent(n, m, v+int64(i)); err != nil {
			return nil, err
		}
		ee = append(ee, e)
	}

	return ee, nil
}

//func DecodeEvent[E any]()
