package store

import (
	"os"

	. "github.com/sokool/stream"
	"github.com/sokool/stream/store/sql"
)

func NewEntitiesX[E Entity]() (Entities[E], error) {
	if cdn := os.Getenv("MYSQL_EVENT_STORE"); cdn != "" {
		c, err := sql.NewConnection(cdn, nil)
		if err != nil {
			return nil, err
		}

		m, err := sql.NewTable[E](c)
		if err != nil {
			return nil, err
		}

		return m, nil
	}

	return NewMemoryEntities[E](), nil
}
