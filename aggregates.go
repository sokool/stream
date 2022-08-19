package stream

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Aggregate[R Root[E], E any] struct {
	// Description
	Description string

	// OnCreate creates stream event appender, where business related events chang
	// state of underlying structure
	OnCreate func(string) (R, error)

	// Events
	Events []E

	// OnSession called on command dispatch when Session not exists
	//OnSession func(R) (Session, error)

	// OnRead when all events are committed to R and state is rebuild from all previously persisted Event
	OnRead RootFunc[R, E]

	// OnRecall
	//OnRecall func(Session, R) error

	// OnWrite called just before events are persisted to database
	OnWrite RootFunc[R, E]

	// OnChange
	OnChange RootFunc[R, E]

	// OnCacheCleanup when aggregate is removed from memory
	OnCacheCleanup RootFunc[R, E]

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
	store EventStore[E]
}

func (a *Aggregate[R, E]) Execute(id string, c RootFunc[R, E]) error {
	for {
		r, err := a.Read(id)
		if err != nil {
			return err
		}

		if err = c(r); err != nil {
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

func (a *Aggregate[R, E]) Read(id string) (r R, err error) {
	if r, err = a.create(id); err != nil {
		return r, err
	}

	var n Namespace
	if n, err = a.namespace(r); err != nil {
		return r, err
	}

	if a.store == nil {
		a.store = NewEventStore[E]()
	}

	es, mm, c := a.store.Stream(n), make([]Event[E], 8), 0
	for {
		switch c, err = es.ReadAt(mm, r.Version()); {

		case err == ErrEndOfStream || err == nil:
			if c == 0 {
				return r, nil
			}

			if failed := a.commit(r, mm[:c]); failed != nil {
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

			return r, a.memory.Set(id, r)

		default:
			return r, Err("%s root read failed due %w", r, err)
		}
	}
}

func (a *Aggregate[R, E]) Write(r R) error {
	mm, err := a.events(r)
	if err != nil {
		return err
	}

	if len(mm) == 0 {
		return nil
	}

	if a.OnWrite != nil {
		if err = a.OnWrite(r); err != nil {
			return err
		}
	}

	var n, m, _ = 0, len(mm), time.Now()
	switch n, err = a.store.Write(mm); {

	case err != nil:
		return err

	case n != m:
		return ErrShortWrite

	default:

	}

	if err = a.commit(r, mm); err != nil {
		return err
	}

	if a.OnChange != nil {
		if err = a.OnChange(r); err != nil {
			return err
		}
	}

	return nil
}

func (a *Aggregate[R, E]) create(id string) (R, error) {
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

func (a *Aggregate[R, E]) namespace(r R) (n Namespace, err error) {
	if n.id, err = NewID(r.ID()); err != nil {
		return n, Err("invalid namespace id %w", err)
	}

	if n.name, err = NewType(r); err != nil {
		return n, Err("invalid namespace name %w", err)
	}

	return n, nil
}

func (a *Aggregate[R, E]) events(r R) (s []Event[E], err error) {
	var n Namespace
	var m Event[E]
	if n, err = a.namespace(r); err != nil {
		return s, nil
	}

	var ee = r.Uncommitted(true)
	for i, e := range ee {
		if m, err = NewEvent(n, e, r.Version()+int64(i)+1); err != nil {
			return s, err
		}
		s = append(s, m)
	}

	return s, nil
}

func (a *Aggregate[R, E]) commit(r R, m []Event[E]) error {
	if len(m) == 0 {
		return nil
	}

	//var s int64
	for i := range m {
		if err := r.Commit(m[i].body, m[i].createdAt); err != nil {
			return err
		}

		//s += int64(len(e.Body))
	}

	//a.version = events[len(events)-1].Sequence
	//a.size += s
	return nil
}

func (a *Aggregate[R, E]) String() (s string) {
	es, ok := a.store.(*eventStore[E])
	if !ok {
		return
	}
	for i := range es.all {
		s += fmt.Sprintf("%s\n", es.all[i])
	}
	return
}

type Root[E any] interface {
	ID() string
	Version() int64
	Uncommitted(clear bool) []E
	Commit(E, time.Time) error
}

type RootFunc[R Root[E], E any] func(R) error

type Context = context.Context

type expired interface{ CacheTimeout() time.Duration }
