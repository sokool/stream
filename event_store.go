package stream

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type EventStore interface {
	ReadWriter(RootID) ReadWriterAt
	Reader(Query) Reader
}

type EventStoreFactory func(Schemas, ...Logger) EventStore

// Query read stream events
type Query struct {
	ID           ID
	Root         Type
	Events       []Type
	FromSequence int64
	From, To     time.Time
	Text         string
	NewestFirst  bool
	Limit        int
	Shutdown     context.Context
}

type store struct {
	mu         sync.Mutex
	namespaces map[RootID]Events
	all        Events
}

func NewEventStore() *store {
	return &store{
		namespaces: map[RootID]Events{},
	}
}

func (s *store) Write(e Events) (n int, err error) {
	for i := range e {
		s.all = append(s.all, e[i])
		s.namespaces[e[i].Root] = append(s.namespaces[e[i].Root], e[i])
	}
	return len(e), nil
}

func (s *store) Reader(q Query) Reader {
	s.mu.Lock()
	defer s.mu.Unlock()

	//sb := Buffer{}
	//st := q.Events
	//
	//for i := range s.all {
	//	if st != nil {
	//		if !st.IsAllowed(s.all[i]) {
	//			continue
	//		}
	//	}
	//	sb.Append(s.all[i])
	//}

	return nil

}

func (s *store) ReadWriter(n RootID) ReadWriterAt {
	return &streamStore{stream: n, store: s}
}

func (s *store) Types() []RootID {
	var st []RootID
	for i := range s.namespaces {
		st = append(st, s.namespaces[i][len(s.namespaces[i])-1].Root)
	}

	return st
}

func (s *store) WriteTo(w Writer) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	nn, err := w.Write(s.all)
	return int64(nn), err
}

func (s *store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.all, s.namespaces = []Event[any]{}, make(map[RootID]Events)
}

func (s *store) Size() (streams int, events int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.namespaces), len(s.all)
}

func (s *store) String() (t string) {
	for i := range s.all {
		t += fmt.Sprintf("%s", &s.all[i])
	}
	return
}

type streamStore struct {
	store  *store
	stream RootID
	mu     sync.Mutex
}

func (s *streamStore) ReadAt(e Events, pos int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	events, ok := s.store.namespaces[s.stream]
	total := len(events)
	max := len(e)

	if !ok || pos > int64(total) {
		return 0, ErrEndOfStream
	}

	var from = int(pos)
	var to = max

	if max >= total {
		to = total
	} else {

		to = from + to
		if to > total {
			to = total
		}

	}

	var i int
	var m Event[any]
	for i, m = range events[from:to] {
		e[i] = m
		//fmt.Println("    ", e)
	}

	//var eos = ""
	//if to >= total {
	//	eos = "EOS"
	//}
	//fmt.Println("read max:", max, "from", from, "to", to, "found", i+1, "total", total, eos)

	if to < total {
		return i + 1, nil
	}

	if int64(total) == pos {
		return 0, ErrEndOfStream
	}

	return i + 1, ErrEndOfStream
}

func (s *streamStore) WriteAt(e Events, pos int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, e := range e {
		if pos >= 0 {
			if int64(len(s.store.namespaces[e.Root])) != pos {
				return i, ErrConcurrentWrite
			}
		}

		s.store.namespaces[e.Root] = append(s.store.namespaces[e.Root], e)
		s.store.all = append(s.store.all, e)
	}

	return len(e), nil
}
