package stream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Root interface {
	ID() string
	Version() int64
	Uncommitted(clear bool) (events []any)
	Commit(event any, createdAt time.Time) error
}

type RootFunc[R Root] func(R) error

type RootID struct {
	id  ID
	typ Type
}

func NewRootID(r Root) (id RootID, err error) {
	if id.id, err = NewID(r.ID()); err != nil {
		return id, Err("invalid namespace id %w", err)
	}

	if id.typ, err = NewType(r); err != nil {
		return id, Err("invalid namespace name %w", err)
	}

	return id, nil
}

func ParseRootID(s string) (n RootID, err error) {
	var p []string
	if p = strings.Split(s, "."); len(p) != 2 {
		return n, Err("wrong `%s` format, please use <id>.<type> ie `N8hY13fsd.Chat`", s)

	}

	if n.id, err = NewID(p[0]); err != nil {
		return n, Err("invalid namespace id %w", err)
	}

	if n.typ, err = NewType(p[1]); err != nil {
		return n, Err("invalid namespace name %w", err)
	}

	return
}

func (id RootID) ID() ID {
	return id.id
}

func (id RootID) Type() Type {
	return id.typ
}

func (id RootID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.String())
}

func (id *RootID) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, null) {
		return nil
	}

	v, err := ParseRootID(string(b))
	if err != nil {
		return err
	}

	*id = v
	return nil
}

func (id RootID) String() string {
	return fmt.Sprintf("%s.%s", id.id, id.typ)
}

func (id RootID) IsZero() bool {
	return id.id == "" || id.typ == ""
}

type Name struct {
	event Type
	root  RootID
}

func NewName(r Root, e Event[any]) (n Name, err error) {
	if n.root, err = NewRootID(r); err != nil {
		return
	}

	return
}
