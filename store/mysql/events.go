package mysql

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	. "github.com/sokool/stream"
	"strings"
	"time"
)

type events struct {
	db  *sqlx.DB
	log Printer
}

func NewEventStore2(host string, p Printer) (*events, error) {
	var s events
	var err error

	db, err := sql.Open("mysql", host)
	if err != nil {
		panic(err)
	}
	s.db = sqlx.NewDb(db, "mysql")

	if err = s.initialise(false); err != nil {
		return nil, err
	}

	if p == nil {
		s.log = func(s string, i ...interface{}) {}
	} else {
		s.log = p
	}

	go func() {
		for range time.NewTicker(time.Second * 10).C {
			if err = s.db.Ping(); err != nil {
				s.log("storage ping failed due %s", err)
			}
		}
	}()

	return &s, nil
}

func (r *events) Stream(n Namespace) ReadWriterAt {
	//TODO implement me
	panic("implement me")
}

func (r *events) Read(q Query) Reader {
	panic("implement me")
}

func (r *events) Write(e Events) (n int, err error) {
	panic("implement me")
}

func (r *events) initialise(drop bool) error {
	q := []string{
		eventsCreate,
	}

	if drop {
		q = append([]string{eventsDrop}, q...)
	}

	return r.execute(func(tx transaction) error {
		for i := range q {
			if _, err := tx.Exec(q[i]); err != nil && !strings.Contains(err.Error(), "Duplicate key name") {
				return err
			}
		}

		return nil
	})
}

func (r *events) execute(fn func(transaction) error) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			tx.Rollback()

		} else {
			// all good, commit
			if err = tx.Commit(); err != nil {
				//fmt.Println("transact.commit", err)
			}
		}
	}()

	err = fn(tx)

	return err
}

type transaction interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

const (
	eventsCreate = `
CREATE TABLE IF NOT EXISTS events (
    stream          VARCHAR(255) NOT NULL,
    type            VARCHAR(255),
    name            VARCHAR(255) NOT NULL,
    sequence        INT,
    author          VARCHAR(255),
    created_at      TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP(6),
    body            JSON,
    meta            JSON,
    id              VARCHAR(36) NOT NULL,
    correlation_id  VARCHAR(36),
    causation_id    VARCHAR(36),

    PRIMARY KEY (type, stream, sequence),
	UNIQUE(id),
    INDEX(created_at ASC)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_unicode_ci;`

	eventsDrop = `DROP TABLE IF EXISTS aggregates_events;`

	eventsIndexA = `CREATE INDEX stream_created_at ON aggregates_events (stream, created_at);`
	eventsIndexB = `CREATE INDEX type_created_at ON aggregates_events (type, created_at);`
	eventsIndexC = `CREATE INDEX created_at_type ON aggregates_events (created_at, type);`
	eventsIndexD = `CREATE INDEX author_created_at ON aggregates_events (author, created_at);`
	eventsIndexE = `CREATE INDEX name_created_at ON aggregates_events (name, created_at);`
	eventsIndexF = `CREATE INDEX name ON aggregates_events (name);`

	eventsWrite = `
INSERT INTO
    aggregates_events(
        stream, type, name, sequence, author, created_at, body, meta, id, 
        correlation_id, causation_id) 
    VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	eventsRead = `
SELECT * 
    FROM aggregates_events
    WHERE
        type = :type           AND
        stream = :stream       AND
        sequence > :sequence 
    ORDER BY
        created_at,
        sequence,
        stream ASC`

	eventsLastSequence = `
SELECT IFNULL(max(sequence), 0) as sequence 
	FROM aggregates_events 
	WHERE 
		type = ? AND 
		stream = ?`
)

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
//
//func (s *stream) Write(events []Event[any]) (n int, err error) {
//	return s.WriteAt(events, -1)
//}

type parameters = map[string]any

type reader struct {
	db *sqlx.DB

	ctx  context.Context
	from int64

	select_
}

func NewReader(db *events) *reader {
	return &reader{
		db:      db.db,
		ctx:     context.Background(),
		select_: newEventsQuery("aggregates_events"),
	}
}

func (r *reader) Stream(id string) *reader { r.where("stream", id); return r }

func (r *reader) Name(s string) *reader { r.where("name", s); return r }

func (r *reader) Type(s string) *reader { r.where("type", s); return r }

func (r *reader) Sequence(n int64) *reader { r.from = n; return r }

func (r *reader) Shutdown(c Context) *reader { r.ctx = c; return r }

func (r *reader) ReadAt(e Events, pos int64) (n int, err error) {
	return r.Sequence(pos).Read(e)
}

func (r *reader) Read(e Events) (n int, err error) {
	var size = len(e)

	res, err := r.run(r.ctx, r.db)
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

func newEventsQuery(table string) select_ {
	return select_{
		q: `
SELECT 
    JSON_OBJECT(
               'Stream', stream,
               'Type', type,
               'Name', name,
               'Sequence', sequence,
               'Author', author,
               'CreatedAt', created_at,
               'Body', body,
               'Meta', meta,
               'ID', id,
               'CorrelationID', correlation_id,
               'CausationID', causation_id
           ) json
FROM ` + table + ` 
WHERE %s 
ORDER BY created_at ASC`,
		p: parameters{},
	}
}

type select_ struct {
	q string
	w []string
	p parameters
}

func (s *select_) where(name, value string) *select_ {
	s.w, s.p[name] = append(s.w, fmt.Sprintf("%s = :%s", name, name)), value
	return s
}

func (s *select_) between(name string) {
	//	s.w, s.p[], par["min"] = append(whr, "sequence > :from AND sequence <= :min"), r.from, r.from+int64(size)
	//}
}

func (s *select_) run(c context.Context, db *sqlx.DB) (*sqlx.Rows, error) {
	q := fmt.Sprintf(s.q, strings.Join(s.w, " AND "))
	fmt.Printf("%s\n%v\n", q, s.p)
	return db.NamedQueryContext(c, q, s.p)
}
