package sql

import (
	"strings"

	"github.com/doug-martin/goqu"
	. "github.com/sokool/stream"
)

type EventsReader struct {
	*Connection
	query EventsQuery
}

func NewEventsReader(c *Connection, q Query) *EventsReader {
	return &EventsReader{c, EventsQuery{q, 100}}
}

func (r *EventsReader) ReadAt(e Events, pos int64) (n int, err error) {
	r.query.FromSequence = pos
	return r.Read(e)
}

func (r *EventsReader) Read(e Events) (n int, err error) {
	res, err := r.db.Query(r.query.Limit(e.Size()).String())
	if err != nil {
		return 0, err
	}

	defer res.Close()

	for i := range e {
		if !res.Next() {
			if i == 0 || i < r.query.limit {
				err = ErrEndOfStream
				return
			}

			return
		}

		var b []byte
		if err = res.Scan(&b); err != nil {
			return
		}

		if err = e[n].Decode(b); err != nil {
			return
		}
		n++
	}

	return
}

type EventsQuery struct {
	Query
	limit int
}

func (e EventsQuery) String() string {
	o := func(name string) goqu.OrderedExpression {
		if e.NewestFirst {
			return goqu.I(name).Desc()
		}

		return goqu.I(name).Asc()
	}

	q := goqu.From("aggregates").Select("body").Order(o("created_at"))
	if d := e.Stream; !d.IsEmpty() {
		q = q.Where(goqu.Ex{"id": d.String()}).Order(o("sequence"))
	}

	if t := e.Root; t != "" {
		q = q.Where(goqu.Ex{"root": t})
	}

	if s := e.Events; len(s) != 0 {
		q = q.Where(goqu.Ex{"event": s})
	}

	if n := int(e.FromSequence); n > 0 {
		q = q.Where(goqu.Ex{"sequence": goqu.Op{"gt": n}}).Where(goqu.Ex{"sequence": goqu.Op{"lte": n + e.limit}})
	}

	if s := e.Text; s != "" {
		q = q.Where(goqu.Ex{"body": goqu.Op{"like": "%" + s + "%"}})
	}

	if e.limit > 0 {
		q = q.Limit(uint(e.limit))
	}

	s, _, _ := q.ToSql()
	return strings.ReplaceAll(s, `"`, "")
}

func (e EventsQuery) Limit(n int) EventsQuery { e.limit = n; return e }
