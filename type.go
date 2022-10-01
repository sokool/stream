package stream

import (
	"github.com/google/uuid"
	"reflect"
	"strings"
)

type Type string

func NewType(v any) (Type, error) {
	var s string

	t := reflect.TypeOf(v)
	s = t.Name()
	switch t.Kind() {
	case reflect.Pointer:
		s = t.Elem().Name()
	case reflect.String:
		s = v.(string)
	default:

	}

	if s = strings.ReplaceAll(strings.TrimSpace(s), " ", ""); len(s) == 0 {
		return "", Err("name can not be empty")
	}

	return Type(strings.Title(s)), nil
}

func (t Type) Hash() string {
	return uid(t.String()).String()
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

func uid(s string) ID {
	return ID(uuid.NewSHA1(uuid.NameSpaceDNS, []byte(s)).String())
}
