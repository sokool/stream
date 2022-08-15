package stream

import (
	"strings"
)

type Type string

func NewType(s string) (Type, error) {
	if s = strings.ReplaceAll(strings.TrimSpace(s), " ", ""); len(s) == 0 {
		return "", Err("type string can not be empty")
	}
	return Type(strings.Title(s)), nil
}

func (t Type) String() string {
	return string(t)
}
