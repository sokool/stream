package mysql

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	. "github.com/sokool/stream"
	"strings"
)

type EventsWriter struct {
	*Connection
	root RootID
}

func NewEventsWriter(c *Connection, id ...RootID) *EventsWriter {
	w := EventsWriter{Connection: c}
	if len(id) != 0 {
		w.root = id[0]
	}
	return &w
}

func (w *EventsWriter) WriteAt(events Events, pos int64) (n int, err error) {
	if len(events) == 0 {
		return 0, nil
	}

	var tx *sqlx.Tx

	if tx, err = w.db.BeginTxx(context.TODO(), nil); err != nil {
		return 0, err
	}

	defer func() {
		if err == nil {
			return
		}

		if failed := tx.Rollback(); failed != nil {
			w.log("%s WriteAt rollback failed on %s due %s", w.root, err, failed)
		}
	}()

	for _, e := range events {
		var b []byte
		if !w.root.IsZero() && e.Root != w.root {
			return 0, fmt.Errorf("rootid missmatch, required %s", w.root)
		}

		if b, err = w.schemas.Encode(e); err != nil {
			return 0, err
		}

		q := `INSERT INTO aggregates(id, root, event, sequence, author, created_at, body) VALUES(?, ?, ?, ?, ?, ?, ?)`
		_, err = tx.Exec(q,
			e.Root.ID(),
			e.Root.Type(),
			e.Type,
			e.Sequence,
			e.Meta.Author,
			e.CreatedAt.UTC().Format("2006-01-02 15:04:05.000000"),
			b,
		)

		if err != nil {
			if s := err.Error(); strings.Contains(s, "Error 1062") && strings.Contains(s, "PRIMARY") {
				return 0, ErrConcurrentWrite
			}

			return 0, err
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return len(events), nil
}

func (w *EventsWriter) Write(e Events) (n int, err error) {
	return w.WriteAt(e, -1)
}
