package cache

import (
	"time"
)

type Cache[K any, V any] interface {
	// Add store your item to cache
	Add(key K, value V, expiration time.Duration) bool
	// Get load your item from cache
	Get(key K) (V, bool)
	// Update updates the Value of an existing entry in the cache.
	Update(key K, newValue V, expiration time.Duration) bool
	// Exists check your key exists in cache
	Exists(key K) bool
	// Keys return list of keys in cache
	Keys() []K
	// Delete delete specific item from cache base on key
	Delete(key K) bool
}
