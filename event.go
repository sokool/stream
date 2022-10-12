package stream

import (
	"encoding/json"
	"fmt"
	"time"
)

type Event struct {
	id       ID
	typ      Type
	root     RootID
	sequence int64

	// body TODO
	body any

	// meta TODO
	meta Meta

	// createdAt
	createdAt time.Time

	version int
}

func NewEvent(id RootID, v any, sequence int64) (e Event, err error) {
	var t Type
	if t, err = NewType(v); err != nil {
		return e, nil
	}

	if sequence <= 0 {
		return e, Err("invalid event sequence")
	}

	return Event{
		id:        uid(fmt.Sprintf("%s.%s.%d", id, e.typ, sequence)),
		root:      id,
		typ:       t.CutPrefix(id.Type()),
		sequence:  sequence,
		body:      v,
		meta:      Meta{},
		createdAt: time.Now(),
	}, nil
}

func (e *Event) Stream() ID {
	return e.root.id
}

func (e *Event) Root() Type {
	return e.root.typ
}

func (e *Event) Type() Type {
	return e.typ
}

func (e *Event) Sequence() int64 {
	return e.sequence
}

func (e *Event) Name() string {
	return fmt.Sprintf("%s%s", e.root.typ, e.typ)
}

func (e *Event) Body() any {
	return e.body
}

func (e *Event) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Event) Belongs(to RootID) bool {
	return e.root == to
}

func (e *Event) Correlate(with Event) *Event {
	e.meta.Correlation = with.id
	return e
}

func (e *Event) Respond(to Event) *Event {
	e.meta.Correlation, e.meta.Causation = to.meta.Correlation, to.id
	return e
}

func (e *Event) String() string {
	return fmt.Sprintf("%s:%d:%s[%s]", e.root.id, e.sequence, e.root.typ, e.typ)
}

func (e *Event) GoString() string {
	b, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("%T\n%s\n", e.body, b)
}

func (e *Event) IsZero() bool {
	return e.id == ""
}

func (e *Event) Encode() ([]byte, error) {
	return registry.encode(*e)
}

func (e *Event) Decode(b []byte) error {
	return registry.decode(e, b)
}

type Events []Event

func NewEvents(r Root) (ee Events, err error) {
	var id RootID
	var k = r.Version() + 1
	if id, err = NewRootID(r); err != nil {
		return nil, err
	}

	for i, v := range r.Uncommitted(true) {
		var e Event
		if e, err = NewEvent(id, v, k+int64(i)); err != nil {
			return nil, err
		}

		ee = append(ee, e)
	}

	return ee, nil
}

// Unique gives RootID when all events has same RootID
func (r Events) Unique() RootID {
	if len(r) == 0 || !r.hasUnique(r[0].root) {
		return RootID{}
	}

	return r[0].root
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

func (r Events) hasUnique(id RootID) bool {
	for i := range r {
		if id != r[i].root && !r[i].IsZero() {
			return false
		}
	}
	return true
}

func (r Events) Append(e Event) error {
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
			t, d = append(t, r[i].typ.String()), r[i].sequence
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
	var e Event
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

func (r Events) IsZero() bool {
	return r.Size() == 0
}

func (r Events) Last() Event {
	if s := r.Size(); s != 0 {
		return r[s-1:][0]
	}

	return Event{}

}
