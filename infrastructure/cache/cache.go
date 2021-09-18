package cache

import (
	"sync"
	"time"
)

// Cache is interface to get/set cache item
type Cache interface {
	Get(key string) (interface{}, bool)
	SetWithExpiration(key string, value interface{}, expiration time.Time)
}

// New creates a new Cache
func New() Cache {
	return &cache{
		m: &sync.Map{},
	}
}

type cache struct {
	m *sync.Map
}

// Get returns cache value when key exists.
func (c *cache) Get(k string) (interface{}, bool) {
	v, ok := c.m.Load(k)
	if !ok {
		return nil, false
	}

	i := v.(item)

	if i.expired() {
		c.m.Delete(k)
		return nil, false
	}

	return i.v, true
}

// SetWithExpiration sets cache value by key.
func (c *cache) SetWithExpiration(k string, v interface{}, e time.Time) {
	c.m.Store(
		k,
		item{
			expiration: e,
			v:          v,
		},
	)
}

type item struct {
	v          interface{}
	expiration time.Time
}

func (i item) expired() bool {
	return time.Now().After(i.expiration)
}
