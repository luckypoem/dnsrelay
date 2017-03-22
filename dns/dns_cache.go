package dns

import (
	"sync"
	"time"

	"github.com/miekg/dns"
)

type DNSCacheConfig struct {
	Enable   bool     `toml:"enable"`
	MaxCount int64 `toml:"max-count"`
}

type DomainRecord struct {
	Msg    *dns.Msg
	Ttl    time.Duration
	Expire time.Time
	Hit    int
	HitMu  *sync.RWMutex
}

func (record *DomainRecord) Touch() {
	record.HitMu.Lock()
	defer record.HitMu.Unlock()

	record.Hit++
}

func (record *DomainRecord) IsTimeout() bool {
	return time.Now().After(record.Expire)
}

func (record *DomainRecord) GetCapacity() int64 {
	return 1
}

type DnsCache struct {
	cache *Cache
}

func NewDnsCache(MaxCacity int64) *DnsCache {
	cache := &Cache{MaxCacity:MaxCacity}
	return &DnsCache{cache: cache}
}

func (c *DnsCache) Get(key string) (*dns.Msg, bool) {
	value, ok := c.cache.Get(key)

	if !ok {
		return nil, false
	}

	record := value.(*DomainRecord)

	if record.IsTimeout() {
		c.Remove(key)
		return nil, false
	}

	record.Touch()
	return record.Msg, true
}

func (c *DnsCache) Add(key string, msg *dns.Msg, ttl time.Duration) {

	record := &DomainRecord{
		Msg:msg,
		Ttl:ttl,
		Expire:time.Now().Add(ttl * time.Second),
		Hit:0,
		HitMu:new(sync.RWMutex),
	}
	record.Touch()

	c.cache.Add(key, record)
	return
}

func (c *DnsCache) Remove(key string) {
	c.cache.Remove(key)
}

func (c *DnsCache) Size() int64 {
	return c.cache.Size()
}
