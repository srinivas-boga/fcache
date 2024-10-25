package main

import (
	"fmt"
	"sync"
	"time"
)

// CacheItem is a struct that holds the value and the expiration time of the value
type CacheItem[T any] struct {
	value    T
	expireAt time.Time
}

type Cache[K comparable, V any] struct {
	items map[K]CacheItem[V]
	mu    sync.RWMutex
}

// NewCache creates a new cache
func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		items: make(map[K]CacheItem[V]),
	}
}

// Set adds a new item to the cache
func (c *Cache[K, V]) Set(key K, value V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem[V]{
		value:    value,
		expireAt: time.Now().Add(ttl),
	}
}

// Get returns the value of the key if it exists in the cache
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok {
		var zeroValue V
		return zeroValue, false
	}

	if item.expireAt.Before(time.Now()) {
		var zeroValue V
		return zeroValue, false
	}

	return item.value, true
}

// Delete removes an item from the cache
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// clear removes all key-value pairs from the cache
func (c *Cache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[K]CacheItem[V])
}

// check for expired keys and remove them
func (c *Cache[K, V]) checkExpiredKeysAndDelete() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, item := range c.items {
		if item.expireAt.Before(time.Now()) {
			delete(c.items, key)
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

func main() {
	cache := NewCache[string, string]()
	cache.StartEvictionWorker(1 * time.Second)

	cache.Set("key1", "value1", 1*time.Second)
	cache.Set("key2", "value2", 3*time.Second)

	time.Sleep(2 * time.Second)

	value1, ok1 := cache.Get("key1")
	value2, ok2 := cache.Get("key2")

	fmt.Println(value1, ok1) // Output: false
	fmt.Println(value2, ok2) // Output: true

	cache.Clear()

}
