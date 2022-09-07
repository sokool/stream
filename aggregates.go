package stream

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Aggregate is todo :)
type Aggregate[R Root] struct {
	// Description
	Description string

	// OnCreate creates aggregate root (business entity).
	// It is kept in memory
	OnCreate func(string) (R, error)

	// Events
	//Events []E

	// OnSession called on command dispatch when Session not exists
	//OnSession func(R) (Session, error)

	// OnRead when all events are committed to R and state is rebuild from all previously persisted Event
	OnRead RootFunc[R]

	// OnRecall
	//OnRecall func(Session, R) error

	// OnWrite called just before events are persisted to database
	OnWrite RootFunc[R]

	// OnCommit
	OnCommit func(R, []Event[any]) error

	// OnCacheCleanup when aggregate is removed from memory
	OnCacheCleanup RootFunc[R]

	// RecallAfter default
	//RecallAfter time.Duration

	CleanCacheAfter time.Duration

	//todo think about it
	//EventSourced bool

	LoadEventsInChunks int

	// Logger
	//Log Printer

	//schema *Schemas

	// memory keeps created Changelog of Aggregate in order to avoid rebuilding
	// state of each Aggregate everytime when Command is called
	memory *Cache[string, R]
	mu     sync.Mutex

	// store
	store EventStore
}

func (a *Aggregate[R]) Execute(id string, command RootFunc[R]) error {
	for {
		r, err := a.Read(id)
		if err != nil {
			return err
		}

		if err = command(r); err != nil {
			return err
		}

		switch err = a.Write(r); {
		case err == ErrConcurrentWrite:
			continue

		case err != nil:
			return err
		}

		return nil
	}
}

func (a *Aggregate[R]) Read(id string) (r R, err error) {
	if r, err = a.read(id); err != nil {
		return r, err
	}

	var n Namespace

	if n, err = NewNamespace(r); err != nil {
		return r, err
	}

	if a.store == nil {
		a.store = NewEventStore()
	}

	store, events, m := a.store.Stream(n), make([]Event[any], 8), 0
	for {
		switch m, err = store.ReadAt(events, r.Version()); {

		case err == ErrEndOfStream || err == nil:
			if m == 0 {
				return r, nil
			}

			if failed := a.commit(r, events[:m]); failed != nil {
				return r, failed
			}

			if err == nil {
				continue
			}

			if a.OnRead != nil {
				if err = a.OnRead(r); err != nil {
					return r, err
				}
			}

			return r, nil

		default:
			return r, Err("%s root read failed due %w", r, err)
		}
	}
}

func (a *Aggregate[R]) Write(r R) error {
	events, err := NewEvents(r)
	if err != nil {
		return err
	}

	if len(events) == 0 {
		return nil
	}

	var n int
	switch n, err = a.store.Write(events); {
	case err != nil:
		return err

	case n != len(events):
		return ErrShortWrite

	}

	if err = a.commit(r, events); err != nil {
		return err
	}

	if a.OnWrite != nil {
		if err = a.OnWrite(r); err != nil {
			return err
		}
	}

	return nil
}

func (a *Aggregate[R]) read(id string) (R, error) {
	if a.memory == nil {
		a.memory = NewCache[string, R](a.CleanCacheAfter)
	}

	var r, ok = a.memory.Get(id)
	var err error

	if !ok {
		if r, err = a.OnCreate(id); err != nil {
			return r, err
		}
	}

	return r, nil
}

func (a *Aggregate[R]) commit(r R, e []Event[any]) error {
	if len(e) == 0 {
		return nil
	}

	if a.OnCommit != nil {
		if err := a.OnCommit(r, e); err != nil {
			return err
		}
	}

	//var s int64
	for i := range e {
		if err := r.Commit(e[i].body, e[i].createdAt); err != nil {
			return err
		}

		//s += int64(len(e.Body))
	}

	//a.version = events[len(events)-1].Sequence
	//a.size += s
	return a.memory.Set(r.ID(), r)
}

func (a *Aggregate[R]) String() (s string) {
	es, ok := a.store.(*eventStore)
	if !ok {
		return
	}
	for i := range es.all {
		s += fmt.Sprintf("%s\n", es.all[i].String())
	}
	return
}

type Root interface {
	ID() string
	Version() int64
	Uncommitted(clear bool) (events []any)
	Commit(event any, createdAt time.Time) error
}

type RootFunc[R Root] func(R) error

type Context = context.Context

type expired interface{ CacheTimeout() time.Duration }
