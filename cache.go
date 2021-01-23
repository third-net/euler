package euler

import (
	"fmt"
	"time"

	impl "github.com/hashicorp/golang-lru"
	"github.com/miekg/dns"
)

//Cache simple LRU cache
type Cache struct {
	impl *impl.Cache
}

type cacheEntry struct {
	answer   []dns.RR
	deadline uint32
}

//GET get domain cache
func (c *Cache) GET(domain string) ([]dns.RR, bool) {
	if domain == "" {
		return nil, false
	}
	v, ok := c.impl.Get(domain)
	if !ok {
		fmt.Printf("cache miss domain[%s] reason: not found \r\n", domain)
		return nil, false
	}
	entry := v.(*cacheEntry)
	now := time.Now().Unix()
	if now >= int64(entry.deadline) {
		fmt.Printf("cache miss domain[%s] reason: expired \r\n", domain)
		c.impl.Remove(domain)
		return nil, false
	}
	ttl := int64(entry.deadline) - now
	fmt.Printf("domain[%s] now:%d deadline:[%d] ttl:[%d]\r\n", domain, now, entry.deadline, ttl)
	fmt.Printf("cache hit domain[%s] ttl:%d \r\n", domain, ttl)
	for _, a := range entry.answer {
		a.Header().Ttl = uint32(ttl)
	}
	return entry.answer, true
}

//ADD add domain cache
func (c *Cache) ADD(domain string, answer ...dns.RR) {
	if domain == "" {
		return
	}
	if name := dns.Fqdn(domain); name != answer[0].Header().Name {
		answer[0].Header().Name = name
	}
	ttl := answer[0].Header().Ttl
	if ttl < 1 {
		return
	}
	now := uint32(time.Now().Unix())
	deadline := now + ttl
	c.impl.Add(domain, &cacheEntry{
		answer:   answer,
		deadline: deadline,
	})
	fmt.Printf("cache add  domain[%s] ttl:[%d] deadline[%d] \r\n", domain, ttl, deadline)
}

//NewCache new cache
func NewCache(size int) *Cache {
	impl, _ := impl.New(size)
	return &Cache{impl: impl}
}
