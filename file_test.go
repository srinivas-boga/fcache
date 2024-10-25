package fcache

import (
	"os"
	"testing"
	"time"
)

func TestSaveLoadFromFile(t *testing.T) {
	cache := NewCache[string, int]()
	cache.Set("key1", 1, time.Hour)
	cache.Set("key2", 2, time.Hour)
	cache.Set("key3", 3, time.Hour)

	err := cache.SaveToFile("test.gob")
	if err != nil {
		t.Fatal(err)
	}

	cache2 := NewCache[string, int]()
	err = cache2.LoadFromFile("test.gob", time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	v, ok := cache2.Get("key1")
	if !ok {
		t.Fatal("key1 not found")
	}
	if v != 1 {
		t.Fatalf("expected 1, got %d", v)
	}

	v, ok = cache2.Get("key2")
	if !ok {
		t.Fatal("key2 not found")
	}
	if v != 2 {
		t.Fatalf("expected 2, got %d", v)
	}

	v, ok = cache2.Get("key3")
	if !ok {
		t.Fatal("key3 not found")
	}
	if v != 3 {
		t.Fatalf("expected 3, got %d", v)
	}

	os.Remove("test.gob")
}
