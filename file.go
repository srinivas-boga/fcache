package fcache

import (
	"encoding/gob"
	"os"
)

// This package is used to load items from a file and save items of the cache to a file.

// LoadFromFile loads the cache from a file
func (c *Cache[K, V]) LoadFromFile(path string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	return decoder.Decode(&c.Items)
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
