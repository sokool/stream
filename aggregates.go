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
	onCommit OnEvent

	// onRecall func(Session, R) error
	onRecall func(R) time.Time

	onGrant OnEvent

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

func (a *Aggregates[R]) Get(s Session, id string) (*Aggregate[R], error) {
	var ok bool
	var ar *Aggregate[R]

	d, err := a.typ.NewID(id)
	if err != nil {
		return nil, err
	}
	//if err := s.IsGranted(d.Resource()); err != nil {
	//	return nil, Err("forbidden:%w", err)
	//}
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
		r, err := a.Get(s, id)
		if err != nil {
			return err
		}
		if err = r.Run(c); err != nil {
			return err
		}
		switch err = a.Set(s, r); {
		case err == ErrConcurrentWrite:
			continue

		case err != nil:
			return err
		}

		return nil
	}
}

func (a *Aggregates[R]) Set(s Session, r *Aggregate[R]) error {
	var now = time.Now()
	var err error
	var events = r.Events()
	if a.onCommit != nil {
		for i := range events {
			if err = a.onCommit(s, events[i]); err != nil {
				return err
			}
		}
	}

	if err = a.grant(s, r); err != nil {
		return err
	}

	if err = s.IsGranted(r.Events().Resources()...); err != nil {
		return Err("forbidden:%w", err)
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

func (a *Aggregates[R]) OnCommit(fn OnEvent) *Aggregates[R] {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.onCommit = fn
	return a
}

func (a *Aggregates[R]) OnFirst(gf OnEvent) *Aggregates[R] {
	a.onGrant = gf
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

func (a *Aggregates[R]) grant(s Session, r *Aggregate[R]) error {
	if a.onGrant == nil || r.Version() != 0 {
		return nil
	}
	for _, e := range r.Events() {
		if err := a.onGrant(s, e); err != nil {
			return err
		}
	}
	return nil
}

type Context = context.Context

type Date = time.Time

type expired interface{ CacheTimeout() time.Duration }

//type Cmd[R Root] struct {
//	// id of aggregate that command is
//	ID      string
//	Session Session
//	Command func(R) error
//	// list of session IDs that's allowed to have access to aggregate
//	Allow []string
//}

type OnEvent func(s Session, e Event) error
