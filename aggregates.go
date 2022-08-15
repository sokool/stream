package stream

import (
	"fmt"
	"time"
)

type Aggregates[R Root] struct {
	Aggregate[R]
}

func NewAggregates[R Root](a Aggregate[R]) *Aggregates[R] {
	if a.Store == nil {
		a.Store = NewEvents()
	}

	return &Aggregates[R]{a}
}

func (a *Aggregates[R]) Execute(id ID, c Command[R]) error {
	//r, ok := a.list[id]
	//if !ok {
	//	return Err("id %s not found", id)
	//}
	//
	//if a.conf.OnCreate != nil {
	//	if err := a.conf.OnCreate(r); err != nil {
	//		return err
	//	}
	//}

	r, err := a.read(id)
	if err != nil {
		return err
	}

	return c(r)
	//return command(r)
}

func (a *Aggregates[R]) read(id ID) (R, error) {
	r, err := a.OnCreate(id)
	if err != nil {
		return r, err
	}

	//rw := a.Store.Read(id)

	return r, nil
}

func (a *Aggregates[R]) write(r R) error {
	//id, err := NewID(r.ID(), r.Name())
	//if err != nil {
	//	return err
	//}

	ee := r.Uncommitted(true)
	fmt.Println(ee)
	//a.store.Write()
	//a.list[id] = r

	return nil
}

type Aggregate[R Root] struct {
	Type string

	// OnCreate creates stream event appender, where business related events chang
	// state of underlying structure
	OnCreate func(ID) (R, error)

	// Store
	Store Events

	// Description
	Description string

	// Events
	Events []Event

	// OnSession called on command dispatch when Session not exists
	//OnSession func(R) (Session, error)

	// OnCommand when transaction tries to get access to Aggregate Root (Appender)
	OnCommand func(R) error

	// OnRecall
	//OnRecall func(Session, R) error

	// OnPersist called just before events are persisted to database
	OnPersist func(R) error

	// OnChange
	OnChange func(R, []Message) error

	// OnEvict when aggregate is removed from memory
	OnEvict func(ID) error

	// RecallAfter default
	//RecallAfter time.Duration

	//EvictAfter time.Duration

	//todo think about it
	//EventSourced bool

	//LoadEventsInChunks int

	// Logger
	//Log Printer

	//schema *Schemas
}

type Root interface {
	ID() string
	Name() string
	Uncommitted(clear bool) []Event
	Commit(Event, time.Time) error
}

type Event = any

type Command[R Root] func(R) error
