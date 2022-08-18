package stream

import (
	"reflect"
	"strings"
)

type Type string

func NewType(t any) (Type, error) {
	s := reflect.TypeOf(t)
	n := s.Name()
	if s.Kind() == reflect.Pointer {
		n = s.Elem().Name()
	}
	if s.Kind() == reflect.String {
		n = t.(string)
	}
	if n = strings.ReplaceAll(strings.TrimSpace(n), " ", ""); len(n) == 0 {
		return "", Err("name can not be empty")
	}

	return Type(strings.Title(n)), nil
}

func (t Type) String() string {
	return string(t)
}
