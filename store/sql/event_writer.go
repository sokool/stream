package sql

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	. "github.com/sokool/stream"
)

type EventsWriter struct {
	*Connection
	root Sequence
}

func NewEventsWriter(c *Connection, id ...Sequence) *EventsWriter {
	w := EventsWriter{Connection: c}
	if len(id) != 0 {
		w.root = id[0]
	}
	return &w
}

func (w *EventsWriter) WriteAt(ee Events, pos int64) (n int, err error) {
	if len(ee) == 0 {
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

	for _, e := range ee {
		var b []byte
		if !w.root.IsEmpty() && !e.Belongs(w.root.ID()) {
			return 0, fmt.Errorf("rootid missmatch, required %s", w.root)
		}

		if b, err = e.Encode(); err != nil {
			return 0, err
		}

		q := `INSERT INTO aggregates(id, root, event, sequence, author, created_at, body) VALUES(?, ?, ?, ?, ?, ?, ?)`
		_, err = tx.Exec(q,
			e.Stream().Value(),
			e.Stream().Type(),
			e.Type(),
			e.Sequence(),
			"",
			e.CreatedAt().UTC().Format("2006-01-02 15:04:05.000000"),
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

	return len(ee), nil
}

func (w *EventsWriter) Write(e Events) (n int, err error) {
	return w.WriteAt(e, -1)
}
