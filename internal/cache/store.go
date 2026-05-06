package cache

import (
	"sync"
	"time"
)

type Item struct {
	Value string
	Expiration int64
}

type Cache struct {
	store map[string]Item
	mu sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		store: make(map[string]Item),
	}
}

func (c *Cache) Set(key, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}

	c.store[key] = Item{
		Value: value,
		Expiration: exp,
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	item, exists := c.store[key]
	c.mu.RUnlock()

	if !exists {
		return "", false
	}
	
	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		c.Delete(key)
		return "", false
	}
	
	return item.Value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}

func (c* Cache) StartCleanup(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)

			now := time.Now().UnixNano()

			c.mu.Lock()
			for key, value := range c.store {
				if value.Expiration >0 && now > value.Expiration {
					delete(c.store, key)
				}
			}
			c.mu.Unlock()
		}
	}()
}