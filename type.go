package stream

import (
	"reflect"
	"strings"

	"github.com/google/uuid"
)

type Type string

func NewType[T any](v ...T) (Type, error) {
	var s string
	var r reflect.Type

	if len(v) != 0 {
		switch x := any(v[0]).(type) {
		case string:
			r = reflect.TypeOf(x)
		default:
			r = reflect.TypeOf(x)
		}
	} else {
		var t T
		r = reflect.TypeOf(t)
		s = r.Name()
	}

	switch r.Kind() {
	case reflect.Pointer:
		s = r.Elem().Name()
	case reflect.String:
		if s = "string"; len(v) != 0 {
			s = any(v[0]).(string)
		}
	default:
		s = r.Name()
	}

	if n, ok := Type(s).reformat(); ok {
		return n, nil
	}

	return "", Err("type can not be empty string")
}

func MustType[T any](v ...T) Type {
	t, err := NewType[T](v...)
	if err != nil {
		panic(err)
	}
	return t
}

func (t Type) Rename(s string) Type {
	if v, ok := Type(s).reformat(); ok {
		return v
	}
	return t
}

func (t Type) Hash() UUID {
	return NewUUID(t.String())
}

func (t Type) String() string {
	return string(t)
}

func (t Type) IsZero() bool {
	return t == ""
}

func (t Type) CutPrefix(of Type) Type {
	a, b := t.String(), of.String()
	if strings.Index(a, b) == -1 {
		return t
	}

	return Type(strings.Replace(a, b, "", 1))
}

func (t Type) LowerCase() Type {
	return Type(strings.ToLower(string(t)))
}

func (t Type) reformat() (Type, bool) {
	s := string(t)
	if s = strings.ReplaceAll(s, " ", ""); len(s) == 0 {
		return "", false
	}

	return Type(strings.Title(s)), true
}

func uid(s string) string {
	return uuid.NewSHA1(uuid.NameSpaceDNS, []byte(s)).String()
}

type UUID struct{ id uuid.UUID }

func NewUUID(s string) UUID {
	if s == "" {
		return UUID{}
	}
	return UUID{uuid.NewSHA1(uuid.NameSpaceDNS, []byte(s))}
}

func (u UUID) String() string {
	return u.id.String()
}

func (u UUID) Foo() string {
	return u.String()[:8]
}
func (u UUID) IsEmpty() bool {
	return u.id.Version() == 0
}
