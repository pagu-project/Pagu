package cache

import (
	"sync"
	"time"
)

const (
	_defaultCleanUpCheckDuration = 10 * time.Second
)

type BasicCache[K any, V any] struct {
	cache                sync.Map
	cleanUpCheckDuration time.Duration
	opts                 options
}

type basicCacheEntry[V any] struct {
	Value  V
	Expiry time.Time
}

func NewBasic[K any, V any](cleanUpCheckDuration time.Duration, options ...Option) Cache[K, V] {
	defaultCleanUpDuration := _defaultCleanUpCheckDuration

	opts := defaultServerOptions
	for _, opt := range options {
		opt.apply(&opts)
	}

	if cleanUpCheckDuration != 0 {
		defaultCleanUpDuration = cleanUpCheckDuration
	}

	c := &BasicCache[K, V]{
		cache:                sync.Map{},
		cleanUpCheckDuration: defaultCleanUpDuration,
		opts:                 opts,
	}
	go c.cleanupExpiredEntries()
	return c
}

// Add add new time to cache
//
//   - expiration: 0 for disable expire cache
func (c *BasicCache[K, V]) Add(key K, value V, expiration time.Duration) bool {
	var expiry time.Time
	if expiration != 0 {
		expiry = time.Now().Add(expiration)
	}

	entry := basicCacheEntry[V]{Value: value, Expiry: expiry}
	c.cache.Store(key, entry)
	return true
}

func (c *BasicCache[K, V]) Get(key K) (V, bool) {
	var zeroV V // zero Value of type V
	value, ok := c.cache.Load(key)
	if !ok {
		return zeroV, false
	}

	return value.(basicCacheEntry[V]).Value, true
}

func (c *BasicCache[K, V]) Update(key K, newValue V, expiration time.Duration) bool {
	// Check if the key exists in the cache
	value, ok := c.cache.Load(key)
	if !ok {
		return false // Key not found, nothing to update
	}

	// Update the Value
	entry := value.(basicCacheEntry[V])
	entry.Value = newValue

	// Update the expiration time if a new expiration is provided
	if expiration != 0 {
		entry.Expiry = time.Now().Add(expiration)
	}

	// Store the updated entry back in the cache
	c.cache.Store(key, entry)

	return true
}

func (c *BasicCache[K, V]) Exists(key K) bool {
	_, ok := c.cache.Load(key)
	return ok
}

func (c *BasicCache[K, V]) Keys() []K {
	keys := make([]K, 0)
	c.cache.Range(func(key, _ interface{}) bool {
		keys = append(keys, key.(K))
		return true
	})
	return keys
}

func (c *BasicCache[K, V]) Delete(key K) bool {
	c.cache.Delete(key)
	_, ok := c.cache.Load(key)
	return !ok
}

func (c *BasicCache[K, V]) cleanupExpiredEntries() {
	ticker := time.NewTicker(c.cleanUpCheckDuration) // adjust the cleanup frequency as needed

	for range ticker.C {
		c.cache.Range(func(key, value interface{}) bool {
			entry := value.(basicCacheEntry[V])

			// Skip entries with zero expiration time
			if entry.Expiry.IsZero() {
				return true
			}

			if time.Now().After(entry.Expiry) {
				c.cache.Delete(key)
			}
			return true
		})
	}
}
