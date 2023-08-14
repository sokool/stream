package stream

import (
	"encoding/json"
	"fmt"
	"time"
)

type Event struct {
	sequence Sequence
	typ      Type

	// body TODO
	body event

	// meta TODO
	meta Meta

	// createdAt
	createdAt time.Time

	version int
}

func NewEvent(s Sequence, e event) (Event, error) {
	var t Type
	var err error
	if t, err = NewType(e); err != nil {
		return Event{}, nil
	}

	return Event{
		typ:       t.CutPrefix(s.ID().Type()),
		sequence:  s,
		body:      e,
		meta:      Meta{},
		createdAt: time.Now(),
	}, nil
}

func (e *Event) ID() UUID {
	return e.sequence.Hash()
}

func (e *Event) Type() Type {
	return e.typ
}

func (e *Event) Stream() ID {
	return e.sequence.ID()
}

func (e *Event) Sequence() int64 {
	return e.sequence.Number()
}

func (e *Event) Aggregate() Sequence {
	return e.sequence
}

func (e *Event) Role(id, name string) Role {
	return e.sequence.Resource().Role(id, name)
}

func (e *Event) Name() string {
	return fmt.Sprintf("%s%s", e.Stream().Type(), e.typ)
}

func (e *Event) Body() any {
	return e.body
}

func (e *Event) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Event) Belongs(to ID) bool {
	return e.sequence.ID() == to
}

func (e *Event) Correlate(with Event) *Event {
	e.meta.Correlation = with.ID()
	return e
}

func (e *Event) Respond(to Event) *Event {
	e.meta.Correlation, e.meta.Causation = to.meta.Correlation, to.ID()
	return e
}

func (e *Event) String() string {
	if e.sequence.Number() == 0 {
		return fmt.Sprintf("%s[%s]", e.Stream(), e.Type())
	}
	return fmt.Sprintf("%s[%s]#%d", e.Stream(), e.Type(), e.Sequence())
}

func (e *Event) GoString() string {
	b, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("%T\n%s\n", e.body, b)
}

func (e *Event) IsEmpty() bool {
	return e.sequence.IsEmpty()
}

func (e *Event) Encode() ([]byte, error) {
	return registry.encode(*e)
}

func (e *Event) Decode(b []byte) error {
	return registry.decode(e, b)
}

func (e *Event) Resource() Resource {
	r := e.sequence.id
	return Resource{
		ID:     r.uuid.String(),
		Name:   r.typ.String(),
		Action: e.typ.String(),
	}

}

type Events []Event

func NewEvents(s Sequence, events ...event) (ee Events, err error) {
	for i := range events {
		s = s.Next()
		e, err := NewEvent(s, events[i])
		if err != nil {
			return nil, err
		}
		ee = append(ee, e)
	}
	return ee, nil
}

// Unique gives Sequence when all events has same Sequence
func (r Events) Unique() ID {
	if len(r) == 0 || !r.IsUnique(r[0].Stream()) {
		return ID{}
	}

	return r[0].Stream()
}

//func (r Events) Shrink(f Filter) (Events, error) {
//	if f == nil {
//		return r, nil
//	}
//
//	var o Events
//	for i := range r {
//		ok, err := f.Filtrate(&r[i])
//		if err != nil {
//			return nil, err
//		}
//
//		if ok {
//			o = append(o, r[i])
//		}
//	}
//	return o, nil
//}

func (r Events) IsUnique(id ID) bool {
	if id.IsEmpty() {
		return false
	}
	for i := range r {
		if id != r[i].Stream() && !r[i].IsEmpty() {
			return false
		}
	}
	return true
}

func (r Events) String() string {
	var s string
	if id := r.Unique(); !id.IsEmpty() {
		s = id.String()
		var t []string
		var d int64
		for i := range r {
			if r[i].IsEmpty() {
				continue
			}
			t, d = append(t, r[i].typ.String()), r[i].Sequence()
		}

		return fmt.Sprintf("%s%v#%d", id, t, d)
	}

	for i := range r {
		if r[i].IsEmpty() {
			continue
		}
		s += fmt.Sprintf("%s\n", &r[i])
	}

	return s
}

func (r Events) Size() int {
	return len(r)
}

func (r Events) IsEmpty() bool {
	return r.Size() == 0
}

func (r Events) Last() Event {
	if s := r.Size(); s != 0 {
		return r[s-1:][0]
	}

	return Event{}

}

func (r Events) Resources() []Resource {
	var rr []Resource
	for i := range r {
		rr = append(rr, r[i].Resource())
	}
	return rr
}

func (r Events) UUID() UUID {
	if r.IsEmpty() {
		return UUID{}
	}
	return NewUUID(r.String())
}
