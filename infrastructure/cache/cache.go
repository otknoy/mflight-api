//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock_$GOPACKAGE/mock_$GOFILE
package cache

import (
	"sync"
	"time"
)

// Cache is interface to get/set cache item
type Cache interface {
	Get(key string) interface{}
	SetWithExpiration(key string, value interface{}, expiration time.Time)
}

// New creates a new Cache
func New() Cache {
	v := make(map[string]item)

	return &cache{
		v: v,
	}
}

type cache struct {
	mu sync.RWMutex

	v map[string]item
}

// Get returns cache value when key exists.
func (c *cache) Get(k string) interface{} {
	c.mu.RLock()
	i, ok := c.v[k]
	c.mu.RUnlock()

	if !ok {
		return nil
	}

	if i.expired() {
		c.mu.Lock()
		delete(c.v, k)
		c.mu.Unlock()

		return nil
	}

	return i.v
}

// SetWithExpiration sets cache value by key.
func (c *cache) SetWithExpiration(k string, v interface{}, e time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.v[k] = item{
		expiration: e,
		v:          v,
	}
}

type item struct {
	v          interface{}
	expiration time.Time
}

func (i item) expired() bool {
	return time.Now().After(i.expiration)
}
