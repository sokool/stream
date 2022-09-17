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

func (r *Schemas) Apply(e Scheme) error {

	x := append(*r, e)
	r = &x
	return nil
}

func (r *Schemas) Merge(s *Schemas) error {
	return nil
}

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
