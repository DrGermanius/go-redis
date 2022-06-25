package main

import (
	"sync"
	"time"
)

type Cache struct {
	store map[string]Item
	mu    sync.RWMutex
}

type Item struct {
	Value    string
	Deadline time.Time
}

func NewCache() *Cache {
	return &Cache{
		store: make(map[string]Item),
		mu:    sync.RWMutex{},
	}
}

func (c *Cache) Add(k string, v Item) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[k] = v
}

func (c *Cache) Get(k string) (Item, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.store[k]
	return v, ok
}

func (c *Cache) Remove(k string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, k)
}
