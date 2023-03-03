package mysql

import (
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sokool/stream"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Connection struct {
	db  *sqlx.DB
	log stream.Printer
	gdb *gorm.DB
}

func NewConnection(host string, l stream.Printer) (*Connection, error) {
	var c Connection
	var err error

	if c.gdb, err = gorm.Open(mysql.Open(host)); err != nil {
		return nil, err
	}

	if l == nil {
		l = log.Printf
	}

	c.log = l

	c.gdb = c.gdb.Session(&gorm.Session{
		NewDB:  true,
		Logger: logger.New(xx{}, logger.Config{}),
	})

	db, err := c.gdb.DB()
	if err != nil {
		panic(err)
	}

	c.db = sqlx.NewDb(db, "mysql")

	tx, err := c.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			if failed := tx.Rollback(); failed != nil {
				c.log("connection rollback failed due %s", failed)
			}
		}
	}()

	for _, q := range []string{eventsCreate} {
		if _, err = tx.Exec(q); err != nil && !strings.Contains(err.Error(), "Duplicate key name") {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	go func() {
		for range time.NewTicker(time.Second * 10).C {
			if err = c.db.Ping(); err != nil {
				c.log("connection ping failed due %s", err)
			}
		}
	}()

	return &c, nil
}

const (
	eventsCreate = `
CREATE TABLE IF NOT EXISTS aggregates (
    id         	VARCHAR(255) NOT NULL,
    root        VARCHAR(255) NOT NULL,
    event       VARCHAR(255) NOT NULL,
    sequence    INT,
    author      VARCHAR(255),
    created_at  TIMESTAMP(6) NOT NULL,
    body        JSON,

    PRIMARY KEY (id, root, sequence),
    INDEX(created_at ASC)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE utf8mb4_unicode_ci;`

	eventsDrop = `DROP TABLE IF EXISTS aggregates;`

	eventsIndexA = `CREATE INDEX stream_created_at ON aggregates_events (stream, created_at);`
	eventsIndexB = `CREATE INDEX type_created_at ON aggregates_events (type, created_at);`
	eventsIndexC = `CREATE INDEX created_at_type ON aggregates_events (created_at, type);`
	eventsIndexD = `CREATE INDEX author_created_at ON aggregates_events (author, created_at);`
	eventsIndexE = `CREATE INDEX name_created_at ON aggregates_events (name, created_at);`
	eventsIndexF = `CREATE INDEX name ON aggregates_events (name);`
)

type xx struct {
}

func (x xx) Printf(s string, i ...interface{}) {
	//TODO implement me
	panic("implement me")
}
