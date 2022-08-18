package stream

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Events[E any] interface {
	Stream(Namespace) ReadWriterAt[E]
	Read(Query) Reader[E]
	Write(m []Message[E]) (n int, err error)
}

type Query struct {
	Stream     Namespace
	From, To   time.Time
	Descending bool
	Shutdown   context.Context
}

type EventStore[E any] struct {
	mu         sync.Mutex
	namespaces map[Namespace][]Message[E]
	all        []Message[E]
}

func NewEvents[E any]() *EventStore[E] {
	return &EventStore[E]{namespaces: make(map[Namespace][]Message[E])}
}

func (s *EventStore[E]) Write(m []Message[E]) (n int, err error) {
	for i := range m {
		s.all = append(s.all, m[i])
		s.namespaces[m[i].namespace] = append(s.namespaces[m[i].namespace], m[i])
	}
	return len(m), nil
}

func (s *EventStore[E]) Read(q Query) Reader[E] {
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

func (s *EventStore[E]) Stream(n Namespace) ReadWriterAt[E] {
	return &rwp[E]{stream: n, store: s}
}

func (s *EventStore[E]) Types() []Namespace {
	var st []Namespace
	for i := range s.namespaces {
		st = append(st, s.namespaces[i][len(s.namespaces[i])-1].namespace)
	}

	return st
}

func (s *EventStore[E]) WriteTo(w Writer[E]) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	nn, err := w.Write(s.all)
	return int64(nn), err
}

func (s *EventStore[E]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.all, s.namespaces = []Message[E]{}, make(map[Namespace][]Message[E])
}

func (s *EventStore[E]) Size() (streams int, events int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.namespaces), len(s.all)
}

func (s *EventStore[E]) String() (t string) {
	for i := range s.all {
		t += fmt.Sprintf("%s", s.all[i])
	}
	return
}

type rwp[E any] struct {
	store  *EventStore[E]
	stream Namespace
	mu     sync.Mutex
}

func (s *rwp[E]) ReadAt(p []Message[E], pos int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	events, ok := s.store.namespaces[s.stream]
	total := len(events)
	max := len(p)

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
	var m Message[E]
	for i, m = range events[from:to] {
		p[i] = m
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

func (s *rwp[E]) WriteAt(m []Message[E], pos int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, e := range m {
		if pos >= 0 {
			if int64(len(s.store.namespaces[e.namespace])) != pos {
				return i, ErrConcurrentWrite
			}
		}

		s.store.namespaces[e.namespace] = append(s.store.namespaces[e.namespace], e)
		s.store.all = append(s.store.all, e)
	}

	return len(m), nil
}
