package stream

import (
	"sync"
	"time"
)

type Cache[K comparable, V any] struct {
	mu      sync.Mutex
	cleanup func(K, V)
	list    map[K]V
}

func NewCache[K comparable, V any](cleanupAfter ...time.Duration) *Cache[K, V] {
	if len(cleanupAfter) == 0 {
		cleanupAfter = append(cleanupAfter, time.Hour)
	}

	c := Cache[K, V]{
		list: make(map[K]V),
	}

	//go func() {
	//	for range time.NewTimer(cleanupAfter[0]).C {
	//		for range c.list {
	//
	//		}
	//	}
	//}()
	return &c
}

func (c *Cache[K, V]) WithCleanup(fn func(K, V)) *Cache[K, V] {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cleanup = fn
	return c
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	v, ok := c.list[key]
	return v, ok
}

func (c *Cache[K, V]) Set(key K, value V, timeout ...time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.list[key] = value
	return nil
}
