package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	. "github.com/sokool/stream"
)

type EventsStore struct {
	*Connection
}

func NewEventsStore(host string, s *Schemas, l ...Log) (*EventsStore, error) {
	var e EventsStore
	var err error

	if e.Connection, err = NewConnection(host, s, l...); err != nil {
		return nil, err
	}
	e.schemas = s

	return &e, nil
}

func (r *EventsStore) ReadWriter(n RootID) ReadWriterAt {
	return struct {
		*EventsReader
		*EventsWriter
	}{
		NewEventsReader(r.Connection, Query{ID: n.ID(), Root: n.Type()}),
		NewEventsWriter(r.Connection),
	}
}

func (r *EventsStore) Reader(q Query) Reader {
	return NewEventsReader(r.Connection, q)
}

func (r *EventsStore) Write(e Events) (n int, err error) {
	return NewEventsWriter(r.Connection).Write(e)
}
