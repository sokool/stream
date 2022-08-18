package stream

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Namespace struct {
	id   ID
	name Name
}

func NewNamespace(id, name string) (Namespace, error) {
	var n Namespace
	var err error

	if n.id, err = NewID(id); err != nil {
		return n, Err("invalid namespace id %w", err)
	}

	if n.name, err = NewName(name); err != nil {
		return n, Err("invalid namespace name %w", err)
	}

	return n, nil
}

func ParseNamespace(s string) (Namespace, error) {
	if p := strings.Split(s, "."); len(p) == 2 {
		return NewNamespace(p[0], p[1])
	}

	return Namespace{}, Err("wrong id %s string format, please use <value>.<type> ie `N8hY13fsd.Chat`")
}

func MustNamespace(s string) Namespace {
	n, err := ParseNamespace(s)
	if err != nil {
		panic(err)
	}

	return n
}

func (n Namespace) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.String())
}

func (n Namespace) String() string {
	return fmt.Sprintf("%s.%s", n.id, n.name)
}

func (n Namespace) IsZero() bool {
	return n.id == "" || n.name == ""
}
