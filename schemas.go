package stream

import (
	"encoding/json"
	"fmt"
	"github.com/alecthomas/jsonschema"
)

type Schemas []Scheme

// todo instead []byte type RawEvent []byte
func (r *Schemas) Decode(e *Event[any], b []byte) error {
	return json.Unmarshal(b, e)
}

func (r *Schemas) Encode(e Event[any]) ([]byte, error) {
	return json.Marshal(e)
}

func (r *Schemas) Append(e Scheme) error {
	x := append(*r, e)
	r = &x
	return nil
}

func (r *Schemas) Merge(s Schemas) error {
	*r = append(*r, s...)
	return nil
}

func (r *Schemas) Filtrate(e *Event[any]) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Schemas) IsCoupled(with Type, of Events) bool {
	//r.mu.Lock()
	//defer r.mu.Unlock()

	return false
	//for i := range of {
	//	if e := r.get(of[i]); e != nil && e.isCoupled(with) {
	//		return true
	//	}
	//}
	//
	//return false
}

//func (r *Schemas) Names() []string {
//	r.mu.Lock()
//	defer r.mu.Unlock()
//
//	var names []string
//	for n := range r.types {
//		names = append(names, n)
//	}
//
//	for n := range r.aliases {
//		names = append(names, n)
//	}
//
//	return names
//}

//func (r *Schemas) has(e EVENT) (Event, bool) {
//	n := infoOf(e)
//	for _, m := range r.types {
//		if m.info.path == n.path {
//			return *m, true
//		}
//	}
//
//	return Event{}, false
//}

//func (r *Schemas) Get(m Event[any]) *Event {
//	//r.mu.Lock()
//	//defer r.mu.Unlock()
//
//	return r.get(m)
//}

//func (r *Schemas) get(m Event[any]) *Event {
//	n := m.Name()
//	s, ok := r.types[n]
//	if !ok {
//		s = r.aliases[n]
//	}
//
//	return s
//}

type Scheme struct {
	root        RootType
	name        string
	description string
	value       any
	coupling    Coupling
	version     int
	scheme      []byte // describe structure in JSON Event https://json-schema.org/
	info        *info
}

func NewScheme(v any) Scheme {
	s := Scheme{value: v}
	s.info = infoOf(s.value)
	if s.name == "" {
		s.name = s.info.typ
	}

	return s
}

func (e Scheme) Name(s string) Scheme { e.name = s; return e }

func (e Scheme) Description(s string) Scheme { e.description = s; return e }

func (e Scheme) Couple(with ...Type) Scheme {
	e.coupling = e.coupling.Add(with...)
	return e
}

func (e Scheme) String() string {
	return fmt.Sprintf("%s%s", e.root, e.name)
}

func (e Scheme) MarshalJSON() ([]byte, error) {
	return json.Marshal(view{
		"ID":          e.info.uuid,
		"Type":        e.root,
		"Description": e.description,
		"Schema":      jsonschema.Reflect(e.value),
		"Coupling":    e.coupling,
		"Location":    e.info.path,
		"Version":     e.version,
	})
}

func (e Scheme) isCoupled(with ...Type) bool {
	return e.coupling.IsStrong(with...)
}

func (e Scheme) isValid() error {
	if e.root.IsZero() || e.name == "" {
		return fmt.Errorf("schma type, kind, value is required")
	}

	return nil
}

type RootType struct{ Type }
