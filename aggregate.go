package stream

import (
	"fmt"
	"sync"
)

type Aggregate[R Root] struct {
	mu          sync.Mutex
	id          RootID
	version     int64
	uncommitted Events
	root        R
}

func NewAggregate[R Root](id string, nr NewRoot[R], definitions []event) (*Aggregate[R], error) {
	var a Aggregate[R]
	var err error
	if nr == nil {
		return nil, Err("new aggregate requires non-nil NewRoot[R]")
	}

	if a.id, err = NewRootID[R](id); err != nil {
		return nil, Err("new aggregate id failed %w", err)
	}

	if a.root, err = nr(id); err != nil {
		return nil, Err("new aggregate instance failed %w", err)
	}

	if err = registry.set(definitions, a.id.typ); err != nil {
		return nil, Err("new aggregate events definition registration failed %w", err)
	}

	return &a, err
}

func MustAggregate[R Root](id string, nr NewRoot[R], definitions []event) *Aggregate[R] {
	a, err := NewAggregate(id, nr, definitions)
	if err != nil {
		panic(err)
	}
	return a
}

func (a *Aggregate[R]) ID() string {
	return a.id.ID().String()
}

func (a *Aggregate[R]) Name() string {
	return a.id.Type().String()
}

func (a *Aggregate[R]) Version() int64 {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.version
}

func (a *Aggregate[R]) ReadFrom(es EventStore) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	var rw, events, n = es.ReadWriter(a.id), make(Events, 1024), 0
	var err error
	for {
		switch n, err = rw.ReadAt(events, a.version); {

		case err == ErrEndOfStream || err == nil:
			if failed := a.commit(events[:n]); failed != nil {
				return failed
			}

			if err == nil {
				continue
			}

			return nil

		default:
			return Err("%s root read failed due %w", a.id, err)
		}
	}
}

// Run todo recover panic from Command
func (a *Aggregate[R]) Run(c Command[R]) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	var err error
	if a.uncommitted.Size() > 0 {
		return Err("write events before another command in %s aggregate ", a.id)
	}

	if err = c(a.root); err != nil {
		return err
	}

	for i, v := range a.root.Uncommitted(true) {
		e, err := NewEvent(a.id, v, a.version+1+int64(i))
		if err != nil {
			return err
		}
		a.uncommitted = append(a.uncommitted, e)
	}

	return nil
}

func (a *Aggregate[R]) Events() Events {
	return a.uncommitted
}

func (a *Aggregate[R]) WriteTo(es EventStore) (Events, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	var n int
	var err error
	defer func() { a.uncommitted = make(Events, 0) }()
	if a.uncommitted.Size() == 0 {
		return nil, nil
	}

	if ok := registry.exists(a.uncommitted); !ok {
		return nil, Err("%s event schema not found, please register it in stream.Aggregates.Events", a.uncommitted)
	}

	if !a.uncommitted.IsUnique(a.id) {
		return nil, Err("aggregate %s events required to be from one root", a.id)
	}

	switch n, err = es.ReadWriter(a.id).WriteAt(a.uncommitted, a.version); {
	case err != nil:
		return nil, err
	case n != len(a.uncommitted):
		return nil, ErrShortWrite
	default:
		return a.uncommitted, a.commit(a.uncommitted)
	}
}

func (a *Aggregate[R]) String() string {
	a.mu.Lock()
	defer a.mu.Unlock()

	if s := len(a.uncommitted); s != 0 {
		return fmt.Sprintf("%s#%d->%d", a.id, a.version, a.version+int64(len(a.uncommitted)))
	}
	return fmt.Sprintf("%s#%d", a.id, a.version)
}

func (a *Aggregate[R]) commit(e []Event) error {
	if len(e) == 0 {
		return nil
	}

	// todo check if given slice of events has correct iteration of sequences, match it with current
	//      version of R
	for i := range e {
		if err := a.root.Commit(e[i].body, e[i].createdAt); err != nil {
			return err
		}
		a.version++
	}

	return nil
}
