package cache

import (
	"sync"
)

type Cacher interface {
	Size() int
}

type cacheMap map[string]Cacher

type Cache struct {
	data cacheMap
	Size uint64
	mu sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{cacheMap{}, 0, sync.RWMutex{}}
}

func (c *Cache) Add(key string, ca Cacher) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = ca
	c.Size += uint64(ca.Size())
}

func (c *Cache) Get(key string) (Cacher, bool) {
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

