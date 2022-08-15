package stream

import (
	"fmt"
	"strings"
)

// ID ...
type ID struct {
	id  string
	typ Type
}

func NewID(value, typ string) (ID, error) {
	t, err := NewType(typ)
	if err != nil {
		return ID{}, Err("id %w", err)
	}

	if value == "" {
		return ID{}, Err("id is empty")
	}

	return ID{id: value, typ: t}, nil
}

func ParseID(s string) (ID, error) {
	if p := strings.Split(s, "|"); len(p) == 2 {
		return NewID(p[0], p[1])
	}

	if p := strings.Split(s, "."); len(p) == 2 {
		return NewID(p[1], p[0])
	}

	return ID{}, Err("wrong id %s string format, please use <value>.<type> ie `N8hY13fsd.Chat`")
}

func MustID(s string) ID {
	if p := strings.Split(s, "|"); len(p) == 2 {
		id, _ := NewID(p[0], p[1])
		return id
	}

	if p := strings.Split(s, "."); len(p) == 2 {
		id, _ := NewID(p[1], p[0])
		return id
	}

	panic("nope")
}

func (n ID) Value() string {
	return n.id
}

func (n ID) Type() Type {
	return n.typ
}

func (n ID) IsZero() bool {
	return n.id == ""
}

func (n ID) String() string {
	if n.id == "" {
		return n.typ.String()
	}

	return fmt.Sprintf("%s|%s", n.id, n.typ)
}
