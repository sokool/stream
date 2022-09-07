package stream

import (
	"encoding/json"
	"fmt"
)

type Namespace struct {
	stream ID
	root   Type
}

func NewNamespace(r Root) (n Namespace, err error) {
	if n.stream, err = NewID(r.ID()); err != nil {
		return n, Err("invalid namespace id %w", err)
	}

	if n.root, err = NewType(r); err != nil {
		return n, Err("invalid namespace name %w", err)
	}

	return n, nil
}

//func NewNamespace(id, name string) (Namespace, error) {
//	var n Namespace
//	var err error
//
//	if n.id, err = NewID(id); err != nil {
//		return n, Err("invalid namespace id %w", err)
//	}
//
//	if n.name, err = NewType(name); err != nil {
//		return n, Err("invalid namespace name %w", err)
//	}
//
//	return n, nil
//}

//func NewStringNamespace(s string) (n Namespace, err error) {
//	var p []string
//	if p = strings.Split(s, "."); len(p) != 2 {
//		return n, Err("wrong `%s` format, please use <id>.<type> ie `N8hY13fsd.Chat`", s)
//
//	}
//
//	if n.id, err = NewID(p[0]); err != nil {
//		return n, Err("invalid namespace id %w", err)
//	}
//
//	if n.name, err = NewType(p[1]); err != nil {
//		return n, Err("invalid namespace name %w", err)
//	}
//
//	return
//}
//
//func MustNamespace(s string) Namespace {
//	n, err := NewStringNamespace(s)
//	if err != nil {
//		panic(err)
//	}
//
//	return n
//}

func (n Namespace) ID() ID {
	return n.stream
}

func (n Namespace) Type() Type {
	return n.root
}

func (n Namespace) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.String())
}

//func (n *Namespace) UnmarshalJSON(bytes []byte) error {
//
//}

func (n Namespace) String() string {
	return fmt.Sprintf("%s.%s", n.stream, n.root)
}

func (n Namespace) IsZero() bool {
	return n.stream == "" || n.root == ""
}

type Name struct {
	event Type
	root  Namespace
}

func NewName(r Root, e Event[any]) (n Name, err error) {
	if n.root, err = NewNamespace(r); err != nil {
		return
	}

	return
}
