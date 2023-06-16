package sql

import (
	_ "github.com/go-sql-driver/mysql"
	. "github.com/sokool/stream"
)

type EventsStore struct {
	*Connection
}

func NewEventsStore(host string, l Logger) (*EventsStore, error) {
	var e EventsStore
	var err error

	if e.Connection, err = NewConnection(host, l); err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *EventsStore) ReadWriter(s Sequence) ReadWriterAt {
	return struct {
		*EventsReader
		*EventsWriter
	}{
		NewEventsReader(r.Connection, Query{Stream: s.ID().UUID(), Root: s.ID().Type()}),
		NewEventsWriter(r.Connection),
	}
}

func (r *EventsStore) Reader(q Query) Reader {
	return NewEventsReader(r.Connection, q)
}

func (r *EventsStore) Write(e Events) (n int, err error) {
	return NewEventsWriter(r.Connection).Write(e)
}
