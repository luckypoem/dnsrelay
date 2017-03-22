package dns

import (
	"sync"
	"github.com/golang/groupcache/lru"
	"github.com/golang/groupcache"
)


type Capacity interface {
	GetCapacity() int64
}

type CacheItem interface {
	Capacity
}

type Cache struct {
	MaxCacity  int64

	mu         sync.RWMutex
	n          int64

	lru        *lru.Cache
	nhit, nget int64
	nevict     int64 // number of evictions
}

func (c *Cache) Stats() groupcache.CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return groupcache.CacheStats{
		Bytes:     c.n,
		Items:     c.sizeLocked(),
		Gets:      c.nget,
		Hits:      c.nhit,
		Evictions: c.nevict,
	}
}

func (c *Cache) Add(key string, value CacheItem) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = &lru.Cache{
			OnEvicted: func(key lru.Key, value interface{}) {
				val := value.(CacheItem)
				c.n -= val.GetCapacity()
				c.nevict++
			},
		}
	}
	c.lru.Add(key, value)
	c.n += value.GetCapacity()

	for {
		if c.n > c.MaxCacity {
			c.lru.RemoveOldest()
		} else {
			break
		}
	}
}

func (c *Cache) Get(key string) (value CacheItem, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.nget++
	if c.lru == nil {
		return
	}
	vi, ok := c.lru.Get(key)
	if !ok {
		return
	}
	c.nhit++
	return vi.(CacheItem), true
}

func (c *Cache) Remove(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.lru.Remove(key)
}

func (c *Cache) Capacity() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.n
}

func (c *Cache) Size() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.sizeLocked()
}

func (c *Cache) sizeLocked() int64 {
	if c.lru == nil {
		return 0
	}
	return int64(c.lru.Len())
}