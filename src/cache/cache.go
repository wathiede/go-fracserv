package cache

import (
	"sync"
)

type Cacheable interface {
	Size() int
}

type Cache struct {
	data map[string]Cacheable
	Size uint64
	mu sync.RWMutex
}

func (c *Cache) Add(key string, ca Cacheable) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = ca
	c.Size += uint64(ca.Size())
}

func (c *Cache) Get(key string) (Cacheable, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	d, ok := c.data[key]
	return d, ok
}

func (c *Cache) Del(key string) {
	// TODO(wathiede): Implement this and an MRU eviction policy
	/*
	c.mu.Lock()
	defer c.mu.RUnlock()
	ca := c.Data[key]
	c.Size -= ca.Size()
	*/
}

