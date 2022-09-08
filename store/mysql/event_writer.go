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

func (w *EventsWriter) WriteAt(e Events, sequence int64) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func (w *EventsWriter) Write(e Events) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

//func (s *stream) WriteAt(events []Event[any], pos int64) (n int, err error) {
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
//		if pos, err = s.last(events[0].Stream); err != nil {
//			return 0, fmt.Errorf("reading last sequence in a stream failed due %s", err)
//		}
//	}
//
//	if tx, err = s.conn.BeginTxx(s.termination, nil); err != nil {
//		return 0, err
//	}
//
//	defer func() {
//		if err == nil {
//			return
//		}
//
//		if failed := tx.Rollback(); failed != nil {
//			s.log("%s WriteAt rollback failed on %s due %s", s.namespace, err, failed)
//		}
//	}()
//
//	for i, event := range events {
//		sequence.next = pos + int64(i) + 1
//
//		if event.Stream != s.namespace {
//			return 0, stream.ErrWrongName
//		}
//
//		if sequence.auto {
//			event.Sequence = sequence.next
//		}
//
//		if sequence.next != event.Sequence {
//			return 0, stream.ErrWrongSequence
//		}
//
//		if err = s.store(event, tx.Exec); err != nil {
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
//	return len(events), nil
//}
