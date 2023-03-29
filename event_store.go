package stream

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Repository interface {
	EventStore(schemas) (EventStore, error)
}

type EventStore interface {
	ReadWriter(Sequence) ReadWriterAt
	Reader(Query) Reader
}

type EventStoreFunc func(*schemas, Printer) EventStore

// Query read stream events
type Query struct {
	Stream       UUID
	Root         Type
	Events       []Type
	FromSequence int64
	From, To     time.Time
	Text         string
	NewestFirst  bool
	Shutdown     context.Context
}

// todo use https://github.com/hashicorp/go-memdb
type store struct {
	mu         sync.Mutex
	namespaces map[ID]Events
	all        Events
}

var MemoryEventStore = NewMemoryEventStore()

func NewMemoryEventStore() *store {
	return &store{
		namespaces: map[ID]Events{},
	}
}

func (s *store) Write(e Events) (n int, err error) {
	for i := range e {
		s.all = append(s.all, e[i])
		id := e[i].Stream()
		s.namespaces[id] = append(s.namespaces[id], e[i])
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

func (s *store) ReadWriter(n Sequence) ReadWriterAt {
	return &streamStore{stream: n, store: s}
}

func (s *store) Types() []ID {
	var st []ID
	for i := range s.namespaces {
		st = append(st, s.namespaces[i][len(s.namespaces[i])-1].Stream())
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

	s.all, s.namespaces = []Event{}, make(map[ID]Events)
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
	stream Sequence
	mu     sync.Mutex
}

func (s *streamStore) ReadAt(e Events, pos int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	events, ok := s.store.namespaces[s.stream.ID()]
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
	var m Event
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

	id := s.stream.ID()
	if !e.IsUnique(id) {
		return 0, Err("wrong root")
	}

	if int64(s.store.namespaces[id].Size()) != pos {
		return 0, ErrConcurrentWrite
	}

	s.store.namespaces[id] = append(s.store.namespaces[id], e...)
	s.store.all = append(s.store.all, e...)

	return len(e), nil
}
