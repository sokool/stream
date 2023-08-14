package stream

import "fmt"

type Session interface {
	Grant(...Role) error
	IsGranted(...Resource) error
}

type Role struct {
	ID        string
	Name      string
	Resources []Resource
}

func (r Role) Resource(id, name, action string) Role {
	r.Resources = append(r.Resources, Resource{id, name, action})
	return r
}

func (r Role) String() (s string) {
	x := fmt.Sprintf("%s:%s", r.Name, r.ID)
	for i := range r.Resources {
		s += fmt.Sprintf("%s:%s", x, r.Resources[i])
	}
	return s
}

type Resource struct {
	ID     string
	Name   string
	Action string
}

func (r Resource) String() string {
	return fmt.Sprintf("%s:%s:%s", r.Name, r.ID, r.Action)
}

func (r Resource) Role(id, name string) Role {
	return Role{id, name, []Resource{r}}
}
