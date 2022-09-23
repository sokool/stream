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

	// CreatedAt
	CreatedAt time.Time

	Coupled []Type
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
		ID:        uid(fmt.Sprintf("%s.%s.%d", id, e.Type, sequence)),
		Root:      id,
		Type:      t.CutPrefix(id.Type()),
		Sequence:  sequence,
		Body:      v,
		Meta:      Meta{},
		CreatedAt: time.Now(),
	}, nil
}

func (e *Event[T]) Correlate(d ID) *Event[T] {
	e.Meta.Correlation = d
	return e
}

func (e *Event[T]) Respond(to Event[any]) *Event[T] {
	e.Meta.Correlation, e.Meta.Causation = to.Meta.Correlation, to.ID
	return e
}

func (e *Event[E]) String() string {
	return fmt.Sprintf("%s:%d:%s[%s]", e.Root.id, e.Sequence, e.Root.typ, e.Type)
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
	var k = r.Version() + 1
	if id, err = NewRootID(r); err != nil {
		return nil, err
	}

	for i, v := range r.Uncommitted(true) {
		var e Event[any]
		if e, err = NewEvent(id, v, k+int64(i)); err != nil {
			return nil, err
		}

		ee = append(ee, e)
	}

	return ee, nil
}

// Unique gives RootID when all events has same RootID
func (r Events) Unique() RootID {
	if len(r) == 0 || !r.hasUnique(r[0].Root) {
		return RootID{}
	}

	return r[0].Root
}

func (r Events) Shrink(f Filter) (Events, error) {
	if f == nil {
		return r, nil
	}

	var o Events
	for i := range r {
		ok, err := f.Filtrate(&r[i])
		if err != nil {
			return nil, err
		}

		if ok {
			o = append(o, r[i])
		}
	}
	return o, nil
}

func (r Events) hasUnique(id RootID) bool {
	for i := range r {
		if id != r[i].Root && !r[i].IsZero() {
			return false
		}
	}
	return true
}

func (r Events) Append(e Event[any]) error {
	return nil
}

func (r Events) String() string {
	var s string
	if id := r.Unique(); !id.IsZero() {
		s = id.String()
		var t []string
		var d int64
		for i := range r {
			if r[i].IsZero() {
				continue
			}
			t, d = append(t, r[i].Type.String()), r[i].Sequence
		}

		return fmt.Sprintf("%s:%d:%s%v", id.id, d, id.typ, t)
	}

	for i := range r {
		if r[i].IsZero() {
			continue
		}
		s += fmt.Sprintf("%s\n", &r[i])
	}

	return s
}

func (r Events) Extend(s Root) (err error) {
	var id RootID
	var e Event[any]
	var k = s.Version() + 1
	if id, err = NewRootID(s); err != nil {
		return err
	}

	for i, v := range s.Uncommitted(true) {
		if e, err = NewEvent(id, v, k+int64(i)); err != nil {
			return err
		}

		if err = r.Append(e); err != nil {
			return err
		}
	}

	return nil
}

func (r Events) Size() int {
	return len(r)
}

//func DecodeEvent[E any]()
