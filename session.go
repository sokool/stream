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

func (r Role) String() string {
	s := fmt.Sprintf("role:%s:%s", r.Name, r.ID)
	for i := range r.Resources {
		s += fmt.Sprintf("\n%s", r.Resources[i])
	}
	return s
}

type Resource struct {
	ID     string
	Name   string
	Action string
}

func (r Resource) String() string {
	return fmt.Sprintf("resource:%s:%s:%s", r.Name, r.ID, r.Action)
}

func (r Resource) Role(id, name string) Role {
	return Role{id, name, []Resource{r}}
}
