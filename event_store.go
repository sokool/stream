package stream

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type EventStore[E any] interface {
	Stream(Namespace) ReadWriterAt[E]
	Read(Query) Reader[E]
	Write(m []Event[E]) (n int, err error)
}

type Query struct {
	Stream     Namespace
	From, To   time.Time
	Descending bool
	Shutdown   context.Context
}

type eventStore[E any] struct {
	mu         sync.Mutex
	namespaces map[Namespace][]Event[E]
	all        []Event[E]
}

func NewEventStore[E any]() EventStore[E] {
	return &eventStore[E]{namespaces: make(map[Namespace][]Event[E])}
}

func (s *eventStore[E]) Write(m []Event[E]) (n int, err error) {
	for i := range m {
		s.all = append(s.all, m[i])
		s.namespaces[m[i].namespace] = append(s.namespaces[m[i].namespace], m[i])
	}
	return len(m), nil
}

func (s *eventStore[E]) Read(q Query) Reader[E] {
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

func (s *eventStore[E]) Stream(n Namespace) ReadWriterAt[E] {
	return &rwp[E]{stream: n, store: s}
}

func (s *eventStore[E]) Types() []Namespace {
	var st []Namespace
	for i := range s.namespaces {
		st = append(st, s.namespaces[i][len(s.namespaces[i])-1].namespace)
	}

	return st
}

func (s *eventStore[E]) WriteTo(w Writer[E]) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	nn, err := w.Write(s.all)
	return int64(nn), err
}

func (s *eventStore[E]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.all, s.namespaces = []Event[E]{}, make(map[Namespace][]Event[E])
}

func (s *eventStore[E]) Size() (streams int, events int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.namespaces), len(s.all)
}

func (s *eventStore[E]) String() (t string) {
	for i := range s.all {
		t += fmt.Sprintf("%s", s.all[i])
	}
	return
}

type rwp[E any] struct {
	store  *eventStore[E]
	stream Namespace
	mu     sync.Mutex
}

func (s *rwp[E]) ReadAt(e []Event[E], pos int64) (int, error) {
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
	var m Event[E]
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

func (s *rwp[E]) WriteAt(e []Event[E], pos int64) (int, error) {
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
