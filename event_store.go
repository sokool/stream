package stream

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type EventStore interface {
	Stream(Namespace) ReadWriterAt
	Read(Query) Reader
	Write([]Event[any]) (n int, err error)
}

type Query struct {
	Stream     Namespace
	From, To   time.Time
	Descending bool
	Shutdown   context.Context
}

type eventStore struct {
	mu         sync.Mutex
	namespaces map[Namespace][]Event[any]
	all        []Event[any]
}

func NewEventStore() EventStore {
	return &eventStore{namespaces: make(map[Namespace][]Event[any])}
}

func (s *eventStore) Write(e []Event[any]) (n int, err error) {
	for i := range e {
		s.all = append(s.all, e[i])
		s.namespaces[e[i].namespace] = append(s.namespaces[e[i].namespace], e[i])
	}
	return len(e), nil
}

func (s *eventStore) Read(q Query) Reader {
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

func (s *eventStore) Stream(n Namespace) ReadWriterAt {
	return &rwp{stream: n, store: s}
}

func (s *eventStore) Types() []Namespace {
	var st []Namespace
	for i := range s.namespaces {
		st = append(st, s.namespaces[i][len(s.namespaces[i])-1].namespace)
	}

	return st
}

func (s *eventStore) WriteTo(w Writer) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	nn, err := w.Write(s.all)
	return int64(nn), err
}

func (s *eventStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.all, s.namespaces = []Event[any]{}, make(map[Namespace][]Event[any])
}

func (s *eventStore) Size() (streams int, events int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.namespaces), len(s.all)
}

func (s *eventStore) String() (t string) {
	for i := range s.all {
		t += fmt.Sprintf("%s", s.all[i])
	}
	return
}

type rwp struct {
	store  *eventStore
	stream Namespace
	mu     sync.Mutex
}

func (s *rwp) ReadAt(e []Event[any], pos int64) (int, error) {
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

func (s *rwp) WriteAt(e []Event[any], pos int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, e := range e {
		if pos >= 0 {
			if int64(len(s.store.namespaces[e.namespace])) != pos {
				return i, ErrConcurrentWrite
			}
		}

		s.store.namespaces[e.namespace] = append(s.store.namespaces[e.namespace], e)
		s.store.all = append(s.store.all, e)
	}

	return len(e), nil
}
