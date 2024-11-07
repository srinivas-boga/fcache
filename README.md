# fcache - In memory thread safe key:value store to cache huge number of entries

[![build passing](https://github.com/srinivas-boga/fcache/actions/workflows/go.yml/badge.svg?branch=main&event=status)](https://github.com/srinivas-boga/fcache/actions/workflows/go.yml)

### Features
* Thread Safe, Concurrent Go-routines can write and read values to the cache.
* Designed to reduce GC overhead for a large number of entries by a custom HashTable Implementation.
* Cache memory size can be specified when creating the cache.
* Simple, easy to read source code.



### Limitations
* Supports only []byte type for keys and values, need to serialize before storing them in the cache.
* Big Key: Value entries that exceed the buffer size are not stored in the cache.


### Future Work
* Currently supports no expiration and eviction on TTL.
* Loading and Storing the cache from / to a file.
