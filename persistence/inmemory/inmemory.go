package inmemory

import (
	"sync"

	"github.com/buraksekili/getir-task/persistence"
)

// InMemory represents a in-memory database.
type InMemory struct {
	storage map[string]string
	mu      *sync.Mutex
}

// New creates a new in-memory persistence layer.
func New(storage map[string]string) InMemory {
	return InMemory{storage, &sync.Mutex{}}
}

// FetchData fetches a data with given key from in-memory.
func (im *InMemory) FetchData(key string) string {
	im.mu.Lock()
	defer im.mu.Unlock()
	return im.storage[key]
}

// StoreData inserts a new key-value pair into in-memory.
func (im *InMemory) StoreData(kv persistence.InMemoryKeyValue) {
	im.mu.Lock()
	defer im.mu.Unlock()
	im.storage[kv.Key] = kv.Value
}
