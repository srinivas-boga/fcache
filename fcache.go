package fcache

import (
	"sync"
	"time"
)

const EVICTION_CHECK_INTERVAL = 1 * time.Second

// CacheItem is a struct that holds the value and the expiration time of the value
type CacheItem[T any] struct {
	Value    T
	ExpireAt time.Time
}

type Cache[K comparable, V any] struct {
	Items map[K]CacheItem[V]
	mu    sync.RWMutex
}

// NewCache creates a new cache
func NewCache[K comparable, V any]() *Cache[K, V] {
	// call eviction worker before returning the cache
	// to start the eviction worker
	cache := &Cache[K, V]{
		Items: make(map[K]CacheItem[V]),
	}
	cache.StartEvictionWorker(EVICTION_CHECK_INTERVAL)
	return cache
}

// Set adds a new item to the cache
func (c *Cache[K, V]) Set(key K, value V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Items[key] = CacheItem[V]{
		Value:    value,
		ExpireAt: time.Now().Add(ttl),
	}
}

// Get returns the value of the key if it exists in the cache
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.Items[key]
	if !ok {
		var zeroValue V
		return zeroValue, false
	}

	if item.ExpireAt.Before(time.Now()) {
		var zeroValue V
		return zeroValue, false
	}

	return item.Value, true
}

// Delete removes an item from the cache
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.Items, key)
}

// clear removes all key-value pairs from the cache
func (c *Cache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Items = make(map[K]CacheItem[V])
}

// check for expired keys and remove them
func (c *Cache[K, V]) checkExpiredKeysAndDelete() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, item := range c.Items {
		if item.ExpireAt.Before(time.Now()) {
			delete(c.Items, key)
		}
	}
}

// StartEvictionWorker starts a goroutine that checks for expired keys and removes them
func (c *Cache[K, V]) StartEvictionWorker(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			c.checkExpiredKeysAndDelete()
		}
	}()
}
