package fcache

const NO_OF_BUCKETS = 512

const BUCKET_SIZE = 64 * 1024

type Cache struct {
	Buckets [NO_OF_BUCKETS]Bucket
}

// NewCache creates a new cache
func NewCache() *Cache {

	var c Cache
	for i := 0; i < NO_OF_BUCKETS; i++ {
		c.Buckets[i] = *NewBucket(BUCKET_SIZE)
	}

	return &c
}

// Set adds a new item to the cache
func (c *Cache) Set(k, v []byte) error {
	bucket := &c.Buckets[hashFnv(k)%NO_OF_BUCKETS]
	return bucket.Set(k, v)
}

// Get returns the value of the key if it exists in the cache
func (c *Cache) Get(k []byte) ([]byte, error) {
	bucket := &c.Buckets[hashFnv(k)%NO_OF_BUCKETS]
	return bucket.Get(k)
}
