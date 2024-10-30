package fcache

import (
	"errors"
	"hash/fnv"
	"sync"
)

// Bucket is a Ring buffer that holds the data of the cache
type Bucket struct {

	// data is the byte slice that holds the data of the bucket
	data []byte
	// mu is used to synchronize access to the cache
	mu sync.RWMutex

	// idx is the index where next data will be written
	idx uint64

	indexMap map[uint64]uint64
}

// NewBucket creates a new bucket
func NewBucket(size uint64) *Bucket {
	b := &Bucket{}
	b.data = make([]byte, size)
	b.idx = 0
	b.indexMap = make(map[uint64]uint64)
	return b
}

func (b *Bucket) Get(k []byte) ([]byte, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	pos, ok := b.indexMap[hashFnv(k)]
	if !ok {
		ErrKeyNotFound := errors.New("key not found")
		return nil, ErrKeyNotFound
	}

	_, v, err := b.ReadAt(pos)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (b *Bucket) ReadAt(idx uint64) ([]byte, []byte, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if idx >= uint64(len(b.data)) {
		ErrIndexOutOfRange := errors.New("index out of range")
		return nil, nil, ErrIndexOutOfRange
	}

	// read the length of the key and the value
	kLen := uint16(b.data[idx])<<8 | uint16(b.data[idx+1])
	vLen := uint16(b.data[idx+2])<<8 | uint16(b.data[idx+3])

	k := make([]byte, kLen)
	v := make([]byte, vLen)

	idx = idx + 4
	n := 0
	for n < int(kLen) {
		read := copy(k[n:], b.data[idx:])
		n += read
		idx += uint64(read)
		if idx == uint64(len(b.data)) {
			idx = 0
		}
	}

	n = 0
	for n < int(vLen) {
		read := copy(v[n:], b.data[idx:])
		n += read
		idx += uint64(read)
		if idx == uint64(len(b.data)) {
			idx = 0
		}
	}

	return k, v, nil
}

// Set adds a new item to the bucket
func (b *Bucket) Set(k, v []byte) error {

	b.mu.Lock()
	defer b.mu.Unlock()

	if len(k) > (1<<16) || len(v) > (1<<16) {
		ErrKeyOrValueTooLarge := errors.New("key or value too large")
		return ErrKeyOrValueTooLarge
	}

	// kvHeader is a 4 byte slice that holds the length of the key and the value
	kvHeader := make([]byte, 4)
	kvHeader[0] = byte(len(k) >> 8)
	kvHeader[1] = byte(len(k))
	kvHeader[2] = byte(len(v) >> 8)
	kvHeader[3] = byte(len(v))

	// kv is the byte slice that holds the key and the value
	pos := b.idx
	kv := append(kvHeader, k...)
	kv = append(kv, v...)

	if len(kv) > len(b.data) {
		ErrDataTooLarge := errors.New("data too large")
		return ErrDataTooLarge
	}

	// write the key and the value to the bucket, overwriting the oldest data if necessary
	n := 0
	for n < len(kv) {
		written := copy(b.data[pos:], kv[n:])

		// remove the elements in the indexMap that are overwritten which have value of pos to pos+written
		for k, v := range b.indexMap {
			if v >= pos && v < pos+uint64(written) {
				delete(b.indexMap, k)
			}
		}
		n += written
		pos += uint64(written)
		if pos == uint64(len(b.data)) {
			pos = 0
		}
	}
	b.indexMap[hashFnv(k)] = b.idx
	b.idx = pos
	return nil
}

// hashFnv computes the FNV-1a hash of the given byte slice
func hashFnv(data []byte) uint64 {
	h := fnv.New64a()
	h.Write(data)
	return h.Sum64()
}
