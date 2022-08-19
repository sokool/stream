package stream

import (
	"fmt"
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
	return uid(t).String()
}

func (t Type) String() string {
	return string(t)
}

func uid(s fmt.Stringer) ID {
	return ID(uuid.NewSHA1(uuid.NameSpaceDNS, []byte(s.String())).String())
}
