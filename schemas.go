package stream

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"
)

type Serializer interface {
	Encode(e Event) ([]byte, error)
	Decode(e *Event, b []byte) error
}

var registry = &schemas{}

type schemas struct {
	mu   sync.Mutex
	list []schema
}

// todo instead []byte type RawEvent []byte
func (r *schemas) decode(e *Event, b []byte) error {
	var j event
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	e.id, e.typ, e.root, e.sequence, e.meta = j.ID, j.Typ, j.Root, j.Sequence, j.Meta
	e.createdAt, e.version = j.CreatedAt, j.Version

	s := r.get(*e)
	if s.isZero() {
		return fmt.Errorf("scheme %s nooot found", s.name())
	}

	if s.migrate != nil {
		s.migrate(e.version, b)
	}
	v := reflect.New(s.reflect)
	if len(b) > 0 {
		if err := json.Unmarshal(j.Body, v.Interface()); err != nil {
			return Err("scheme %s decode failed due %w", e, err)
		}
	}

	e.body = v.Elem().Interface()

	return nil
}

func (r *schemas) encode(e Event) ([]byte, error) {
	b, err := json.Marshal(e.body)
	if err != nil {
		return nil, err
	}

	return json.Marshal(event{
		ID:        e.id,
		Typ:       e.typ,
		Root:      e.root,
		Sequence:  e.sequence,
		Body:      b,
		Meta:      e.meta,
		CreatedAt: e.createdAt,
	})
}

func (r *schemas) merge(s Schemas, root Type) (err error) {
	for e, a := range s {
		var n Type
		var c Types
		if n, err = NewType(e); err != nil {
			return err
		}

		if !a.Transaction.IsZero() {
			c = Types{a.Transaction: true}
		}

		t := reflect.TypeOf(e)
		p := t.PkgPath() + "/" + t.Name()
		r.list = append(r.list, schema{
			id:          uid(p),
			event:       n.CutPrefix(root),
			root:        root,
			description: a.Description,
			coupling:    c,
			version:     0,
			scheme:      nil,
			reflect:     t,
			path:        p,
			migrate:     a.OnMigrate,
		})
	}
	return nil
}

func (r *schemas) get(e Event) schema {
	for _, s := range r.list {
		if s.name() == e.Name() {
			return s
		}
	}
	return schema{}
}

func (r *schemas) Filtrate(e *Event) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *schemas) coupling(e Events) Types {
	r.mu.Lock()
	defer r.mu.Unlock()

	c := make(Types)
	for i := range e {
		c.merge(r.get(e[i]).coupling)
	}

	return c
}

func (r *schemas) isCoupled(t Type, e Events) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range e {
		if r.get(e[i]).coupling[t] {
			return true
		}
	}

	return false
}

type schema struct {
	id          string
	event       Type
	root        Type
	description string
	coupling    Types
	version     int
	scheme      []byte // describe structure in JSON Event https://json-schema.org/
	path        string
	reflect     reflect.Type
	migrate     func(version int, body []byte)
}

func (s schema) name() string {
	return fmt.Sprintf("%s%s", s.root, s.event)
}

func (s schema) isZero() bool {
	return s.root.IsZero() || s.event.IsZero()
}

//	func (s Scheme) Couple(with ...Type) Scheme {
//		s.IsStrongCoupled = s.IsStrongCoupled.Add(with...)
//		return s
//	}
//func (s Scheme) String() string {
//	return s.name()
//}

//	func (s Scheme) MarshalJSON() ([]byte, error) {
//		return json.Marshal(view{
//			"ID":          s.info.uuid,
//			"Type":        s.root,
//			"Description": s.Description,
//			"Schema":      jsonschema.Reflect(s.Event),
//			"IsStrongCoupled":    s.IsStrongCoupled,
//			"Location":    s.info.path,
//			"Version":     s.version,
//		})
//	}
//
//	func (s Scheme) Root(t Type) Scheme {
//		s.Name, s.root = s.Name.CutPrefix(t), t
//		return s
//	}

//
//func (s Scheme) isCoupled(with ...Type) bool {
//	return s.IsStrongCoupled.IsStrong(with...)
//}

type event struct {
	ID       ID //todo root.id root.type event.type event.sequence
	Typ      Type
	Root     RootID
	Sequence int64

	// body TODO
	Body json.RawMessage

	// meta TODO
	Meta Meta

	// createdAt
	CreatedAt time.Time

	Version int
}

type Schemas map[any]Scheme

type Scheme struct {
	Name        string
	Description string
	Transaction Type
	OnMigrate   func(version int, payload []byte)
}
