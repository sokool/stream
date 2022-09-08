package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sokool/stream"
	"log"
	"strings"
	"time"
)

type Log = func(string, ...any)

type Connection struct {
	db  *sqlx.DB
	log Log
}

func NewConnection(host string, l ...Log) (*Connection, error) {
	var c Connection

	db, err := sql.Open("mysql", host)
	if err != nil {
		panic(err)
	}
	c.db = sqlx.NewDb(db, "mysql")

	if err = c.initialise(); err != nil {
		return nil, err
	}

	if len(l) == 0 {
		l = append(l, log.Printf)
	}

	c.log = l[0]

	go func() {
		for range time.NewTicker(time.Second * 10).C {
			if err = c.db.Ping(); err != nil {
				c.log("storage ping failed due %s", err)
			}
		}
	}()

	return &c, nil
}

func (r *Connection) initialise() error {
	q := []string{
		eventsCreate,
	}

	return r.execute(func(tx tx) error {
		for i := range q {
			if _, err := tx.Exec(q[i]); err != nil && !strings.Contains(err.Error(), "Duplicate key name") {
				return err
			}
		}

		return nil
	})
}

func (r *Connection) execute(fn func(tx) error) (err error) {
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

func (r *Connection) events(q stream.Query) *select_ {
	return newEventsQuery("aggregates_events", q)
}

type tx interface {
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

func newEventsQuery(table string, q stream.Query) *select_ {
	s := &select_{
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
ORDER BY :order
LIMIT :limit`,
		p: parameters{
			"order": "created_at ASC",
		},
	}

	if d := string(q.Root.ID); d != "" {
		s.where("stream", d)
		s.p["order"] = "sequence ASC"
	}

	if t := string(q.Root.Type); t != "" {
		s.where("type", t)
	}

	if len(q.Root.Events) != 0 {
		var a []string
		for i := range q.Root.Events {
			a = append(a, q.Root.Events[i].String())
		}
		s.where("name", a...)
	}

	return s
}

type select_ struct {
	q string
	w []string
	p parameters
}

func (s *select_) where(name string, value ...string) *select_ {
	if len(value) == 1 {
		s.w, s.p[name] = append(s.w, fmt.Sprintf("%s = :%s", name, name)), value[0]
		return s
	}

	s.w = append(s.w, fmt.Sprintf(`%s IN("%s")`, name, strings.Join(value, `","`)))
	return s
}

func (s *select_) limit(n int) *select_ {
	s.p["limit"] = n
	return s
}

func (s *select_) between(name string) {
	//	s.w, s.p[], par["min"] = append(whr, "sequence > :from AND sequence <= :min"), r.from, r.from+int64(size)
	//}
}

func (s *select_) run(c context.Context, db *Connection) (*sqlx.Rows, error) {
	q := fmt.Sprintf(s.q, strings.Join(s.w, " AND "))
	fmt.Printf("%s\n%v\n", q, s.p)

	return db.db.NamedQueryContext(c, q, s.p)
}

type parameters = map[string]any
