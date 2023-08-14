package stream

import (
	"fmt"
	"strings"
)

type ID struct {
	uuid UUID
	typ  Type
}

func NewID[T any](s string) (d ID, err error) {
	if d.uuid = NewUUID(s); d.uuid.IsEmpty() {
		return d, Err("id uuid is empty")
	}
	if d.typ, err = NewType[T](); err != nil {
		return d, Err("creating id type failed %w", err)
	}
	return d, nil
}

func ParseID(s string) (ID, error) {
	var p []string
	var id ID
	if p = strings.Split(s, "."); len(p) < 2 {
		return id, Err("parse id `%s` invalid format, please use <id>.<type> ie `N8hY13fsd.Chat`", s)
	}

	id.uuid, id.typ = NewUUID(p[0]), Type(p[1])

	return id, nil
}

func (id ID) UUID() UUID {
	return id.uuid
}

func (id ID) Type() Type {
	return id.typ
}

func (id ID) Hash() UUID {
	return NewUUID(fmt.Sprintf("%s:%s", id.uuid, id.typ))
}

func (id ID) IsEmpty() bool {
	return id.uuid.IsEmpty()
}

func (id ID) String() string {
	return fmt.Sprintf("%s.%s", id.uuid.Short(), id.typ)
}

func (id ID) URN() string {
	return fmt.Sprintf("%s:%s", id.typ, id.uuid)
}

func (id ID) Sequence() Sequence {
	return Sequence{id: id}
}
