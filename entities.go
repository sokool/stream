package stream

import (
	"encoding/json"
	"fmt"
)

type Entity interface {
	ID() string
	Version() int64
}

type Entities[E Entity] interface {
	One(E) error // todo One(Events) (E, error)
	Read([]E, []byte) error
	Update(...E) error
	Delete(...E) error
}

type NewEntities[E Entity] func() Entities[E]

// todo use https://github.com/hashicorp/go-memdb
type entities[E Entity] struct {
	store map[string][]byte
}

func NewMemoryEntities[E Entity]() Entities[E] {
	return &entities[E]{
		store: make(map[string][]byte),
	}
}

func (r *entities[E]) One(d E) error {
	b, found := r.store[d.ID()]
	if !found {
		return ErrDocumentNotFound
	}

	return json.Unmarshal(b, &d)
}

func (r *entities[E]) Read(ee []E, bytes []byte) error {
	var i int
	for _, body := range r.store {
		if err := json.Unmarshal(body, &ee[i]); err != nil {
			return err
		}
		i++
	}
	return nil
}

func (r *entities[E]) Update(e ...E) error {
	for i := range e {
		b, err := json.Marshal(e[i])
		if err != nil {
			return err
		}

		r.store[e[i].ID()] = b
	}
	return nil
}

func (r *entities[E]) Delete(e ...E) error {
	//TODO implement me
	panic("implement me")
}

func (r *entities[E]) Count() int {
	return len(r.store)
}

func (r *entities[E]) String() string {
	var e E
	var s = fmt.Sprintf("%T\n", e)
	for i := range r.store {
		s += fmt.Sprintf("%s\n", string(r.store[i]))
	}
	return s
}
