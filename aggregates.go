package stream

import (
	"context"
	"sync"
	"time"
)

// Aggregate is todo :)
type Aggregate[R Root] struct {
	// Type ...
	Type Type

	// Description
	Description string

	// OnCreate creates aggregate root (business entity).
	// It is kept in memory
	OnCreate func(string) (R, error)

	// OnSession called on command dispatch when Session not exists
	//OnSession func(R) (Session, error)

	// OnLoad when all events are committed to R and state is rebuild from all previously persisted Event
	OnLoad RootFunc[R]

	// OnRecall
	//OnRecall func(Session, R) error

	// OnSave called just after events are persisted to database
	OnSave RootFunc[R]

	// OnCommit when new events are committed to a Root
	OnCommit func(R, Events) error // todo not able to deny it (error)

	// OnRecall
	OnRecall func(R) time.Time

	// OnCacheCleanup when aggregate is removed from memory
	OnCacheCleanup RootFunc[R]

	// Events
	Events Schemas

	CleanCacheAfter time.Duration

	//todo think about it
	//EventSourced bool

	LoadEventsInChunks int

	// Store
	Store EventStore

	// Writer
	Writer Writer

	// Logger
	Log Printer

	// memory keeps created Changelog of Aggregate in order to avoid rebuilding
	// state of each Aggregate everytime when Thread is called
	memory *Cache[string, R]
	mu     sync.Mutex
}

// todo recover panic
func (a *Aggregate[R]) Execute(id string, command RootFunc[R]) error {
	for {
		r, err := a.Get(id)
		if err != nil {
			return err
		}

		if err = command(r); err != nil {
			return err
		}

		switch err = a.Set(r); {
		case err == ErrConcurrentWrite:
			continue

		case err != nil:
			return err
		}

		return nil
	}
}

func (a *Aggregate[R]) Get(id string) (R, error) {
	var found bool
	var d RootID
	var r R
	var err error

	if err = a.init(); err != nil {
		return r, err
	}

	if r, found = a.memory.Get(id); !found {
		if r, err = a.OnCreate(id); err != nil {
			return r, err
		}
	}

	if d, err = NewRootID(r); err != nil {
		return r, err
	}

	rw, evs, m := a.Store.ReadWriter(d), make(Events, a.LoadEventsInChunks), 0
	for {
		switch m, err = rw.ReadAt(evs, r.Version()); {

		case err == ErrEndOfStream || err == nil:
			if failed := a.commit(r, evs[:m]); failed != nil {
				return r, failed
			}

			if err == nil {
				continue
			}

			if a.OnLoad != nil {
				if err = a.OnLoad(r); err != nil {
					return r, err
				}
			}

			return r, a.memory.Set(r.ID(), r)

		default:
			return r, Err("%s root read failed due %w", r, err)
		}
	}
}

func (a *Aggregate[R]) Set(r R) error {
	if err := a.init(); err != nil {
		return err
	}

	events, err := make(Events, 0).Append(r)
	if err != nil || len(events) == 0 {
		return err
	}

	rid := events.Unique()
	if rid.IsZero() { //todo error description
		return Err("aggregate %s events required to be from one root", rid)
	}

	var n int
	switch n, err = a.Store.ReadWriter(rid).WriteAt(events, r.Version()); {
	case err != nil:
		return err

	case n != len(events):
		return ErrShortWrite

	}

	if err = a.commit(r, events); err != nil {
		return err
	}

	if a.Writer != nil {
		if _, err = a.Writer.Write(events); err != nil {
			return err
		}
	}

	if a.OnSave != nil {
		if err = a.OnSave(r); err != nil {
			return err
		}
	}

	return a.memory.Set(r.ID(), r)
}

func (a *Aggregate[R]) commit(r R, e []Event) error {
	if len(e) == 0 {
		return nil
	}

	// todo check if given slice of events has correct iteration of sequences, match it with current
	//      version of R
	for i := range e {
		if err := r.Commit(e[i].body, e[i].createdAt); err != nil {
			return err
		}

		//s += int64(len(e.Body))
	}

	if a.OnCommit != nil {
		if err := a.OnCommit(r, e); err != nil {
			return err
		}
	}
	//a.version = events[len(events)-1].Sequence
	//a.size += s
	return nil
}

func (a *Aggregate[R]) init() (err error) {
	if a.Type.IsZero() {
		var r R
		if a.Type, err = NewType(r); err != nil {
			return
		}
	}

	if a.Store == nil {
		a.Store = NewEventStore()
	}

	if a.memory == nil {
		a.memory = NewCache[string, R](a.CleanCacheAfter)
	}

	if a.LoadEventsInChunks <= 0 {
		a.LoadEventsInChunks = 1024
	}

	return nil
}

func (a *Aggregate[R]) String() string {
	if a.Store == nil {
		return ""
	}

	e := make(Events, 10)
	if _, err := a.Store.Reader(Query{Root: a.Type}).Read(e); err != nil {
		return err.Error()
	}

	return e.String()
}

func (a *Aggregate[R]) Register(in *Service) (err error) {
	if err = a.init(); err != nil {
		return err
	}

	if _, ok := a.Store.(*store); ok {
		a.Store = in.store
	}

	if err = registry.merge(a.Events, a.Type); err != nil {
		return err
	}

	if a.Log == nil {
		a.Log = in.logger(a.Type)
	}

	if a.Writer != nil {
		if err = in.register(a.Writer, a.Type); err != nil {
			return err
		}
	}
	a.Writer = in
	return nil
}

type Context = context.Context

type Date = time.Time

type expired interface{ CacheTimeout() time.Duration }
