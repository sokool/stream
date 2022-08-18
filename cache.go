package stream

import (
	"time"
)

type Cache[K comparable, V any] struct {
	list map[K]V
}

func NewCache[K comparable, V any](cleanupAfter ...time.Duration) *Cache[K, V] {
	if len(cleanupAfter) == 0 {
		cleanupAfter = append(cleanupAfter, time.Hour)
	}

	c := Cache[K, V]{
		list: make(map[K]V),
	}

	go func() {
		for range time.NewTimer(cleanupAfter[0]).C {
			for range c.list {

			}
		}
	}()
	return &c
}

func (c Cache[K, V]) Get(key K) (V, bool) {

	v, ok := c.list[key]
	return v, ok
}

func (c Cache[K, V]) Set(key K, value V, timeout ...time.Duration) error {
	c.list[key] = value
	return nil
}
