package fcache

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestBucketGetSet(t *testing.T) {
	b := NewBucket(100)
	key := []byte("key")
	val := []byte("value")
	err := b.Set(key, val)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	v, err := b.Get(key)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !bytes.Equal(val, v) {
		t.Errorf("expected %v, got %v", val, v)
	}
}

func TestBucketGetSetOverwritten(t *testing.T) {

	b := NewBucket(100)
	key1 := []byte("key1")
	val1 := make([]byte, 30)
	rand.Read(val1)
	key2 := []byte("key2")
	val2 := make([]byte, 50)
	rand.Read(val2)
	key3 := []byte("key3")
	val3 := make([]byte, 30)
	rand.Read(val3)

	// failing need to fix the overwritten case in the bucket

	err := b.Set(key1, val1)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	err = b.Set(key2, val2)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	err = b.Set(key3, val3)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// expect overwritten key1, found key2 and key3
	_, err = b.Get(key1)
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	v, err := b.Get(key2)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !bytes.Equal(val2, v) {
		t.Errorf("expected %v, got %v", val2, v)
	}

	v, err = b.Get(key3)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !bytes.Equal(val3, v) {
		t.Errorf("expected %v, got %v", val3, v)
	}

}
