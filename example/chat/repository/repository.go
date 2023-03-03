package repository

import (
	"os"

	"github.com/sokool/stream"
	"github.com/sokool/stream/store/mysql"
)

func storage[E stream.Entity](sne stream.NewEntity[E]) (stream.Entities[E], error) {
	if cdn := os.Getenv("MYSQL_EVENT_STORE"); cdn != "" {
		c, err := mysql.NewConnection(cdn, nil)
		if err != nil {
			return nil, err
		}

		m, err := mysql.NewTable[E](c, sne)
		if err != nil {
			return nil, err
		}

		return m, nil
	}

	return stream.NewEntities[E](sne), nil
}
