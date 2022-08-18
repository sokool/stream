package stream

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Events interface {
	Stream(Namespace) ReadWriterAt
	Read(Query) Reader
	Write(m []Message) (n int, err error)
}

type Query struct {
	Stream     Namespace
	From, To   time.Time
	Descending bool
	Shutdown   context.Context
}

type EventStore struct {
	mu         sync.Mutex
	namespaces NamespacedMessages
	all        []Message
}

type NamespacedMessages map[Namespace][]Message

func NewEvents() *EventStore {
	return &EventStore{namespaces: make(NamespacedMessages)}
}

func (s *EventStore) Write(m []Message) (n int, err error) {
	for i := range m {
		s.all = append(s.all, m[i])
		s.namespaces[m[i].sequence.namespace] = append(s.namespaces[m[i].sequence.namespace], m[i])
	}
	return len(m), nil
}

func (s *EventStore) Read(q Query) Reader {
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

func (s *EventStore) Stream(n Namespace) ReadWriterAt {
	return &rwAt{stream: n, store: s}
}

//func (s *EventStore) ReadFrom(r Reader) (int64, error) {
//	m, err := Read(r)
//	if err != nil && err != EOS {
//		return 0, err
//	}
//	s.all = append(s.all, m...)
//	return int64(len(m)), nil
//}

func (s *EventStore) Types() []Namespace {
	var st []Namespace
	for i := range s.namespaces {
		st = append(st, s.namespaces[i][len(s.namespaces[i])-1].sequence.namespace)
	}

	return st
}

func (s *EventStore) WriteTo(w Writer) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	nn, err := w.Write(s.all)
	return int64(nn), err
}

func (s *EventStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.all, s.namespaces = []Message{}, make(map[Namespace][]Message)
}

func (s *EventStore) Size() (streams int, events int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.namespaces), len(s.all)
}

func (s *EventStore) String() (t string) {
	for i := range s.all {
		t += fmt.Sprintf("%s", s.all[i])
	}
	return
}

func (s *EventStore) readAt(p []Message, n Namespace, at int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	events, ok := s.namespaces[n]
	total := len(events)
	max := len(p)

	if !ok || at > int64(total) {
		return 0, ErrEndOfStream
	}

	var from = int(at)
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
	var event Message
	for i, event = range events[from:to] {
		p[i] = event
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

	if int64(total) == at {
		return 0, ErrEndOfStream
	}

	return i + 1, ErrEndOfStream
}

func (s *EventStore) writeAt(m Message, at int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if at >= 0 {
		if int64(len(s.namespaces[m.sequence.namespace])) != at {
			return ErrConcurrentWrite
		}
	}

	s.namespaces[m.sequence.namespace] = append(s.namespaces[m.sequence.namespace], m)
	s.all = append(s.all, m)

	return nil
}

type rwAt struct {
	store  *EventStore
	stream Namespace
}

func (s *rwAt) ReadAt(events []Message, pos int64) (n int, err error) {
	return s.store.readAt(events, s.stream, pos)
}

func (s *rwAt) WriteAt(events []Message, pos int64) (n int, err error) {
	for i := range events {
		if err = s.store.writeAt(events[i], pos+int64(i)); err != nil {
			return i, err
		}
	}

	return len(events), nil
}
