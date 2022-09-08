package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	. "github.com/sokool/stream"
)

type EventsStore struct {
	*Connection
}

func NewEventsStore(host string, l ...Log) (*EventsStore, error) {
	var s EventsStore
	var err error

	if s.Connection, err = NewConnection(host, l...); err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *EventsStore) ReadWriter(n RootID) ReadWriterAt {
	var q Query
	q.Root.ID, q.Root.Type = n.ID(), n.Type()

	return struct {
		*EventsReader
		*EventsWriter
	}{
		NewEventsReader(r.Connection, q),
		NewEventsWriter(r.Connection),
	}
}

func (r *EventsStore) Reader(q Query) Reader {
	return NewEventsReader(r.Connection, q)
}

func (r *EventsStore) Write(e Events) (n int, err error) {
	panic("implement me")
}

//type stream struct {
//	*events
//	namespace   Namespace
//	termination Context
//}
//
//func (s *stream) ReadAt(events []Event[any], pos int64) (n int, err error) {
//	return s.read(s.termination, events, `
//SELECT *
//   FROM aggregates_events
//   WHERE
//       type = :type           AND
//       stream = :stream       AND
//       sequence > :from AND sequence <= :min
//
//   ORDER BY created_at ASC`, parameters{
//		"stream": s.namespace.ID(),
//		"type":   s.namespace.Type(),
//		"from":   pos,
//		"min":    pos + int64(len(events)),
//	})
//
//}
//
