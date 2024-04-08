package main

import (
	"sync"

	"cache.example/lru"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.LRU
	cacheBytes int64
}

func (c *cache) Add(key string, val ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, val)
}

func (c *cache) get(key string) (b ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	v, ok := c.lru.Get(key)
	if !ok {
		return
	}
	b = v.(ByteView)
	return b, true
}
