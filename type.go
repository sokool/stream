package stream

import (
	"strings"
)

type Name string

func NewName(s string) (Name, error) {
	if s = strings.ReplaceAll(strings.TrimSpace(s), " ", ""); len(s) == 0 {
		return "", Err("name can not be empty")
	}
	return Name(strings.Title(s)), nil
}

func (t Name) String() string {
	return string(t)
}
