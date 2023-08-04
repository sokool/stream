package stream

import (
	"context"
	"sync"
	"time"
)

type NewRoot[R Root] func(string) (R, error)

// Aggregates is todo :)
type Aggregates[R Root] struct {
	typ         Type
	description string

	// onCreate creates aggregate root (business entity).
	// It is kept in memory
	onCreate NewRoot[R]

	// OnSession called on command dispatch when Session not exists
	//OnSession func(R) (Session, error)

	// onLoad when all events are committed to R and state is rebuild from all previously persisted Event
	onLoad Command[R]

	// onSave called just after events are persisted to database
	onSave Command[R]

	// onCommit when new events are committed to a Root
	onCommit func(R, Events) error // todo not able to deny it (error)

	// onRecall func(Session, R) error
	onRecall func(R) time.Time

	// onCacheCleanup when aggregate is removed from memory
	onCacheCleanup Command[R]

	definitions []event

	// Events
	//Events Schemas

	loadEventsInChunks int

	// store
	store EventStore

	// writer
	writer Writer

	// Logger
	log Logger

	// memory keeps created Changelog of Aggregates in order to avoid rebuilding
	// state of each Aggregates everytime when Thread is called
	memory *Cache[ID, *Aggregate[R]]
	mu     sync.Mutex
}

func NewAggregates[R Root](rf NewRoot[R], definitions []event) *Aggregates[R] {
	rt := MustType[R]()
	return &Aggregates[R]{
		onCreate:           rf,
		typ:                rt,
		definitions:        definitions,
		store:              MemoryEventStore,
		memory:             NewCache[ID, *Aggregate[R]](time.Minute),
		log:                newLogger(rt.String()),
		loadEventsInChunks: 1024,
	}

}

func (a *Aggregates[R]) Get(id string) (*Aggregate[R], error) {
	var ok bool
	var ar *Aggregate[R]

	d, err := a.typ.NewID(id)
	if err != nil {
		return nil, err
	}
	if ar, ok = a.memory.Get(d); !ok {
		if ar, err = NewAggregate[R](d, a.onCreate, a.definitions); err != nil {
			return nil, err
		}
	}

	if err = ar.ReadFrom(a.store); err != nil {
		return nil, err
	}

	if a.onLoad != nil {
		if err = a.onLoad(ar.root); err != nil {
			return nil, err
		}
	}

	return ar, a.memory.Set(ar.sequence.id, ar)
}

func (a *Aggregates[R]) Execute(s Session, id string, c Command[R]) error {
	for {
		r, err := a.Get(id)
		if err != nil {
			return err
		}
		if err := s.IsGranted(r.String()); err != nil {
			return Err("forbidden:%w", err)
		}
		if err = r.Run(c); err != nil {
			return err
		}
		if err = s.IsGranted(r.Events().String()); err != nil {
			return Err("forbidden:%w", err)
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

func (a *Aggregates[R]) Set(r *Aggregate[R]) error {
	var now = time.Now()
	var err error
	var events Events
	if a.onCommit != nil {
		if err = a.onCommit(r.root, events); err != nil {
			return err
		}
	}

	if events, err = r.WriteTo(a.store); err != nil {
		return err
	}

	if events.Size() == 0 {
		return nil
	}

	a.log("dbg: %s stored in %s", events, time.Since(now))

	if a.writer != nil {
		if _, err = a.writer.Write(events); err != nil {
			return err
		}
	}

	if a.onSave != nil {
		if err = a.onSave(r.root); err != nil {
			return err
		}
	}
	return a.memory.Set(r.sequence.id, r)
}

func (a *Aggregates[R]) String() string {
	e := make(Events, 10)
	if _, err := a.store.Reader(Query{Root: a.typ}).Read(e); err != nil {
		return err.Error()
	}

	return e.String()
}

func (a *Aggregates[R]) Compose(e *Engine) error {
	a.
		Storage(e.store).
		Logger(e.logger).
		Writer(e)
	return nil
}

func (a *Aggregates[R]) CacheInterval(d time.Duration) *Aggregates[R] {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.memory = NewCache[ID, *Aggregate[R]](d)
	return a
}

func (a *Aggregates[R]) Rename(s string) *Aggregates[R] {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.typ = a.typ.Rename(s)
	a.log = newLogger(a.typ.String())
	return a
}

func (a *Aggregates[R]) OnCommit(fn func(R, Events) error) *Aggregates[R] {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.onCommit = fn
	return a
}

func (a *Aggregates[R]) Storage(es EventStore) *Aggregates[R] {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.store = es
	return a
}

func (a *Aggregates[R]) Writer(w Writer) *Aggregates[R] {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.writer = w
	return a
}

func (a *Aggregates[R]) Logger(l NewLogger) *Aggregates[R] {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.log = l(a.typ.String())
	return a
}

type Context = context.Context

type Date = time.Time

type expired interface{ CacheTimeout() time.Duration }

type Session interface {
	IsGranted(resource string) error
}
