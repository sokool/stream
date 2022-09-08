package mysql

import (
	. "github.com/sokool/stream"
)

type EventsWriter struct {
	db *Connection
}

func NewEventsWriter(c *Connection) *EventsWriter {
	return &EventsWriter{c}
}

func (w *EventsWriter) WriteAt(events Events, pos int64) (n int, err error) {
	//	if len(events) == 0 {
	//		return 0, nil
	//	}
	//
	//	var tx *sqlx.Tx
	//	var sequence struct {
	//		auto bool
	//		next int64
	//	}
	//
	//	// auto sequence when pos is -1, it means that each sequence of events is
	//	// overwritten.
	//	if sequence.auto = pos == -1; sequence.auto {
	//		if pos, err = s.last(events[0].Namespace().ID().String()); err != nil {
	//			return 0, fmt.Errorf("reading last sequence in a stream failed due %s", err)
	//		}
	//	}
	//
	//	if tx, err = w.db.db.BeginTxx(context.TODO(), nil); err != nil {
	//		return 0, err
	//	}
	//
	//	defer func() {
	//		if err == nil {
	//			return
	//		}
	//
	//		if failed := tx.Rollback(); failed != nil {
	//			//s.log("%s WriteAt rollback failed on %s due %s", s.namespace, err, failed)
	//		}
	//	}()
	//
	//	for i, e := range events {
	//		sequence.next = pos + int64(i) + 1
	//
	//		if e.Namespace() != w.namespace {
	//			return 0, fmt.Errorf("rootid missmatch")
	//		}
	//
	//		if sequence.auto {
	//			e.Sequence = sequence.next
	//		}
	//
	//		if sequence.next != e.Sequence() {
	//			return 0, fmt.Errorf("wrong sequence of event")
	//		}
	//
	//		if err = w.store(e, tx.Exec); err != nil {
	//			if s := err.Error(); strings.Contains(s, "Error 1062") && strings.Contains(s, "PRIMARY") {
	//				return 0, stream.ErrConcurrentWrite
	//			}
	//
	//			return 0, err
	//		}
	//	}
	//
	//	if err = tx.Commit(); err != nil {
	//		return 0, err
	//	}
	//
	return len(events), nil
}

func (w *EventsWriter) Write(e Events) (n int, err error) {
	return w.WriteAt(e, -1)
}
