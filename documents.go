package stream

import (
	"fmt"
	"reflect"
)

type Entity interface {
	ID() string
	Version() int64
}

type EntityFunc[E Entity] func(Events) (E, error)

type Documents[E Entity] interface {
	Create(Events) (E, error) // todo One(Events) (E, error)
	Read([]byte) ([]E, error)
	Update(...E) error
	Delete(...E) error
}

// todo use https://github.com/hashicorp/go-memdb
type documents[E Entity] struct {
	create EntityFunc[E]
	store  map[string]E
}

func NewDocuments[E Entity](fn EntityFunc[E]) Documents[E] {
	return &documents[E]{
		store:  make(map[string]E),
		create: fn,
	}
}

func (r *documents[E]) Create(e Events) (E, error) {
	d, err := r.create(e)
	if err != nil || reflect.ValueOf(d).IsNil() {
		return d, err
	}

	b, found := r.store[d.ID()]
	if !found {
		return d, ErrDocumentNotSupported
	}

	return b, nil
}

func (r *documents[E]) Read(bytes []byte) (ee []E, _ error) {
	for _, body := range r.store {
		ee = append(ee, body)
	}
	return ee, nil
}

func (r *documents[E]) Update(e ...E) error {
	for i := range e {
		r.store[e[i].ID()] = e[i]
	}
	return nil
}

func (r *documents[E]) Delete(e ...E) error {
	//TODO implement me
	panic("implement me")
}

func (r *documents[E]) Count() int {
	return len(r.store)
}

func (r *documents[E]) String() string {
	var e E
	var s = fmt.Sprintf("%T\n", e)
	for i := range r.store {
		s += fmt.Sprintf("%v\n", r.store[i])
	}
	return s
}
