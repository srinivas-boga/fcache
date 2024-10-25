package fcache

import (
	"encoding/gob"
	"os"
	"time"
)

// This package is used to load items from a file and save items of the cache to a file.

// LoadFromFile loads the cache from a file
func (c *Cache[K, V]) LoadFromFile(path string, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&c.Items); err != nil {
		return err
	}

	// set the expiration time for the items as the current time with a TTL of 1 hour
	for key, item := range c.Items {
		c.Items[key] = CacheItem[V]{
			Value:    item.Value,
			ExpireAt: time.Now().Add(ttl),
		}
	}
	return nil
}

// SaveToFile saves the cache to a file
func (c *Cache[K, V]) SaveToFile(path string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(c.Items)
}
