package stream

import "fmt"

type Session interface {
	Grant(Role, ...Resource) error
	IsGranted(...Resource) error
}

type Role struct {
	ID   string
	Name string
}

type Resource struct {
	ID     string
	Name   string
	Action string
}

func (r Resource) String() string {
	return fmt.Sprintf("%s:%s:%s", r.Name, r.ID, r.Action)
}
