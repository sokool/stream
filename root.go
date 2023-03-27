package stream

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Root interface {
	//Entity
	Committer
	Uncommitted(clear bool) []event
}

type Committer interface {
	Commit(e event, createdAt time.Time) error //  todo not able to deny it (remove error)
}

type Command[R Root] func(R) error

// RootID
// todo
//   - rename to ID
//   - consider attributes as string, type, sequence
type RootID struct {
	id  ID
	typ Type
}

func NewRootID[T any](id string) (RootID, error) {
	var rid RootID
	var err error
	var t T
	if rid.id, err = NewID(id); err != nil {
		return rid, Err("invalid id string %w", err)
	}
	if rid.typ, err = NewType(t); err != nil {
		return rid, Err("invalid type name %w", err)
	}
	return rid, nil
}

func MustID[T any](id string) RootID {
	rid, err := NewRootID[T](id)
	if err != nil {
		panic(err)
	}
	return rid
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

func (id RootID) Hash() string {
	return uid(id.String())
}

func (id RootID) MarshalJSON() ([]byte, error) {
	return json.Marshal(view{
		"ID":   id.id.String(),
		"Type": id.typ.String(),
	})
}

func (id *RootID) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, null) {
		return nil
	}

	var v view
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	r, err := ParseRootID(fmt.Sprintf("%s.%s", v["ID"], v["Type"]))
	if err != nil {
		return err
	}

	*id = r
	return nil
}

func (id RootID) String() string {
	if id.IsZero() {
		return ""
	}
	return fmt.Sprintf("%s:%s", string(id.id), id.typ)
}

func (id RootID) IsZero() bool {
	return id.id == "" || id.typ == ""
}
