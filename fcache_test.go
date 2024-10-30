package fcache

import (
	"bytes"
	"testing"
)

func TestCacheSetGet(t *testing.T) {
	cache := NewCache()
	key := []byte("key")
	val := []byte("value")

	err := cache.Set(key, val)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	v, err := cache.Get(key)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if !bytes.Equal(val, v) {
		t.Errorf("expected %v, got %v", val, v)
	}
}
