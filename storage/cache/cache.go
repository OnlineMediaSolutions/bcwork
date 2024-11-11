package cache

import (
	"sync"
)

const (
	// Cache keys
	HistoryNewValueCacheKey = "history_new_value"
	HistoryOldValueCacheKey = "history_old_value"
)

type Cache interface {
	Get(key string) (any, bool)
	Set(key string, value any)
	Delete(key string)
}

type InMemoryCache struct {
	cache *sync.Map
}

var _ Cache = (*InMemoryCache)(nil)

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		cache: &sync.Map{},
	}
}

func (c *InMemoryCache) Get(key string) (any, bool) {
	return c.cache.Load(key)
}

func (c *InMemoryCache) Set(key string, value any) {
	c.cache.Store(key, value)
}

func (c *InMemoryCache) Delete(key string) {
	c.cache.Delete(key)
}
