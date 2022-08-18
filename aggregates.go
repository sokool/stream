package stream

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Aggregate[R Root] struct {
	Name string

	// Description
	Description string

	// OnCreate creates stream event appender, where business related events chang
	// state of underlying structure
	OnCreate func(ID) (R, error)

	// Events
	Events []Event

	// OnSession called on command dispatch when Session not exists
	//OnSession func(R) (Session, error)

	// OnRead when all events are committed to R and state is rebuild from all previously persisted Event
	OnRead func(Sequence, R) error

	// OnRecall
	//OnRecall func(Session, R) error

	// OnWrite called just before events are persisted to database
	OnWrite func(Sequence, R) error

	// OnChange
	OnChange func(Sequence, R, []Message) error

	// OnCacheCleanup when aggregate is removed from memory
	OnCacheCleanup func(Sequence, R) error

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
	memory *Cache[Namespace, Rootx[R]]
	mu     sync.Mutex

	// store
	store Events
}

func (a *Aggregate[R]) Execute(id string, command func(R) error) error {
	for {
		r, err := a.read(id)
		if err != nil {
			return err
		}

		if err = command(r.root); err != nil {
			return err
		}

		switch err = a.write(r); {
		case err == ErrConcurrentWrite:
			continue

		case err != nil:
			return err
		}

		return nil
	}
}

func (a *Aggregate[R]) Read(id string) (R, error) {
	r, err := a.read(id)
	if err != nil {
		return r.root, err
	}

	return r.root, nil
}

func (a *Aggregate[R]) Write(r R) error {
	return nil
}

func (a *Aggregate[R]) create(n Namespace) (Rootx[R], error) {
	if a.memory == nil {
		a.memory = NewCache[Namespace, Rootx[R]](a.CleanCacheAfter)
	}

	var r, ok = a.memory.Get(n)
	var err error

	if !ok {
		if r.root, err = a.OnCreate(n.id); err != nil {
			return r, err
		}

		if r.sequence, err = NewSequence(n); err != nil {
			return r, err
		}
	}

	return r, nil
}

func (a *Aggregate[R]) read(id string) (Rootx[R], error) {
	var r Rootx[R]
	var n Namespace
	var err error

	if n, err = NewNamespace(id, a.Name); err != nil {
		return r, err
	}

	if r, err = a.create(n); err != nil {
		return r, err
	}

	if a.store == nil {
		a.store = NewEvents()
	}

	var (
		es     = a.store.Stream(n)
		events = make([]Message, 8)
		m      int
	)

	for {
		switch m, err = es.ReadAt(events, r.sequence.number); {

		case err == ErrEndOfStream || err == nil:
			if m == 0 {
				return r, nil
			}

			if failed := a.commit(r, events[:m]); failed != nil {
				return r, failed
			}

			if r.sequence = events[m-1].sequence; err == nil {
				continue
			}

			if a.OnRead != nil {
				if err = a.OnRead(r.sequence, r.root); err != nil {
					return r, err
				}
			}

			return r, a.memory.Set(r.sequence.namespace, r)

		case err != nil:
			return r, Err("%s root read failed due %w", r.root, err)
		}
	}
}

func (a *Aggregate[R]) set(r Rootx[R]) error {
	//if d := c.evict(); d > 0 {
	//	r.memory.Set(c.stream.String(), c, d)
	//	return nil
	//}
	//
	//r.memory.Delete(c.stream.String())

	a.memory.Set(r.sequence.namespace, r)

	return nil
}

func (a *Aggregate[R]) write(r Rootx[R]) error {
	var err error

	var events []Message
	for _, e := range r.root.Uncommitted(true) {
		r.sequence.number++
		events = append(events, NewMessage(r.sequence, e))
	}

	if len(events) == 0 {
		return nil
	}

	//if ctx := s.Context(); ctx != nil {
	//	select {
	//	case <-ctx.Done():
	//		return nil, ctx.Err()
	//	default:
	//
	//	}
	//}

	//events, err := a.aggregate.schema.Convert(s, a.stream, a.version, une...)
	//if err != nil {
	//	return nil, err
	//}

	if a.OnWrite != nil {
		if err = a.OnWrite(r.sequence, r.root); err != nil {
			return err
		}
	}

	var n, m, _ = 0, len(events), time.Now()
	switch n, err = a.store.Write(events); {

	case err != nil:
		return err

	case n != m:
		return ErrShortWrite

	default:

	}

	if err = a.commit(r, events); err != nil {
		return err
	}

	if a.OnChange != nil {
		if err = a.OnChange(r.sequence, r.root, events); err != nil {
			return err
		}
	}

	//a.Printf("DBG %s committed in %s", MString(events), time.Since(d))

	return nil
}

func (a *Aggregate[R]) commit(r Rootx[R], m []Message) error {
	if len(m) == 0 {
		return nil
	}

	//var s int64
	for i := range m {
		if err := r.root.Commit(m[i].body, m[i].createdAt); err != nil {
			return err
		}

		//s += int64(len(e.Body))
	}

	//a.version = events[len(events)-1].Sequence
	//a.size += s
	return nil
}

func (a *Aggregate[R]) String() (s string) {
	es, ok := a.store.(*EventStore)
	if !ok {
		return
	}
	for i := range es.all {
		s += fmt.Sprintf("%s\n", es.all[i])
	}
	return
}

type Aggregates map[Name]Aggregate[Root]

type Rootx[R Root] struct {
	root     R
	sequence Sequence
}

type Root interface {
	Uncommitted(clear bool) []Event
	Commit(Event, time.Time) error
}

type Event = any

type Context = context.Context

//type Command[R Root] func(R) error
//type NewRoot[R Root] func(ID) (R, error)
