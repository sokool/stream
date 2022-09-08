package mysql

import (
	"context"
	. "github.com/sokool/stream"
)

type EventsReader struct {
	db *Connection

	ctx  context.Context
	from int64
	*select_
}

func NewEventsReader(c *Connection, q Query) *EventsReader {
	return &EventsReader{
		db:      c,
		ctx:     context.Background(),
		select_: c.events(q),
	}
}

func (r *EventsReader) ReadAt(e Events, pos int64) (n int, err error) {
	return r.Read(e)
}

func (r *EventsReader) Read(e Events) (n int, err error) {
	var size = len(e)

	res, err := r.limit(size).run(r.ctx, r.db)
	if err != nil {
		return 0, err
	}

	defer res.Close()

	for n = range e {
		if !res.Next() {
			if n == 0 || n < size {
				return n, ErrEndOfStream
			}

			return n, nil
		}

		var b []byte
		if err = res.Scan(&b); err != nil {
			return
		}

		if err = e[n].UnmarshalJSON(b); err != nil {
			return
		}
	}

	return
}
