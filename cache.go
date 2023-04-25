package stream

import (
	"time"

	"github.com/Code-Hex/go-generics-cache"
)

type Cache[K comparable, V any] struct {
	engine *cache.Cache[K, V]
}

func NewCache[K comparable, V any](d ...time.Duration) *Cache[K, V] {
	if len(d) == 0 {
		d = append(d, time.Hour)
	}

	c := Cache[K, V]{
		cache.New[K, V](cache.WithJanitorInterval[K, V](d[0])),
	}

	return &c
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	return c.engine.Get(key)
}

func (c *Cache[K, V]) Set(key K, value V, timeout ...time.Duration) error {
	if len(timeout) == 0 {
		timeout = append(timeout, time.Minute*10)
	}
	c.engine.Set(key, value, cache.WithExpiration(timeout[0]))
	return nil
}
