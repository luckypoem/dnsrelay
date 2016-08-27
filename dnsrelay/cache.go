package dnsrelay

import (
	"sync"
	"time"

	"github.com/miekg/dns"
)

type DNSCache struct {
	Backend  string
	Expire   int
	Maxcount int
}

type KeyNotFound struct {
	key string
}

func (e KeyNotFound) Error() string {
	return e.key + " " + "not found"
}

type KeyExpired struct {
	Key string
}

func (e KeyExpired) Error() string {
	return e.Key + " " + "expired"
}

type CacheIsFull struct{}
func (e CacheIsFull) Error() string {
	return "Cache is Full"
}

type SerializerError struct{}

func (e SerializerError) Error() string {
	return "Serializer error"
}

type DomainRecord struct {
	Msg    *dns.Msg
	Ttl    time.Duration
	Expire time.Time
	Hit    int
	mu     *sync.RWMutex

}

func (record *DomainRecord) Touch() {
	record.mu.Lock()
	defer record.mu.Unlock()

	record.Hit++
}

func (record *DomainRecord) IsTimeout() bool {
	record.mu.RLock()
	defer record.mu.RUnlock()

	return time.Now().After(record.Expire)
}


type Cache interface {
	Get(key string) (Msg *dns.Msg, err error)
	Set(key string, Msg *dns.Msg) error
	Exists(key string) bool
	Remove(key string)
	Length() int
	Serve()
}

type MemoryCache struct {
	Backend    map[string]DomainRecord
	DefaultTtl time.Duration
	Maxcount   int
	mu         sync.RWMutex
}

func (c *MemoryCache) Get(key string) (*dns.Msg, error) {
	c.mu.RLock()
	record, ok := c.Backend[key]
	c.mu.RUnlock()

	if !ok {
		return nil, KeyNotFound{key}
	}

	if record.IsTimeout() {
		c.Remove(key)
		return nil, KeyExpired{key}
	}

	record.Touch()
	return record.Msg, nil
}

func (c *MemoryCache) Set(key string, msg *dns.Msg, ttl time.Duration) error {
	if c.Full() && !c.Exists(key) {
		return CacheIsFull{}
	}

	if ttl == 0 {
		ttl = c.DefaultTtl
	}

	record := DomainRecord{
		Msg:msg,
		Ttl:ttl,
		Expire:time.Now().Add(ttl * time.Second),
		Hit:0,
		mu:new(sync.RWMutex),
	}
	record.Touch()

	c.mu.Lock()
	c.Backend[key] = record
	c.mu.Unlock()
	return nil
}

func (c *MemoryCache) Remove(key string) {
	c.mu.Lock()
	delete(c.Backend, key)
	c.mu.Unlock()
}

func (c *MemoryCache) Exists(key string) bool {
	c.mu.RLock()
	_, ok := c.Backend[key]
	c.mu.RUnlock()
	return ok
}

func (c *MemoryCache) Length() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.Backend)
}

func (c *MemoryCache) Full() bool {
	if c.Maxcount == 0 {
		return false
	}
	return c.Length() >= c.Maxcount
}

func (c *MemoryCache) Serve() {
	go func() {
		tick := time.Tick(30 * time.Second)
		for _ = range tick {
			c.mu.Lock()
			for domain, record := range c.Backend {
				if record.IsTimeout() {
					delete(c.Backend, domain)
				}
			}
			c.mu.Unlock()
		}
	}()
}

