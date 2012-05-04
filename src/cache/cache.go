// Copyright 2012 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
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
	mu   sync.RWMutex
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
