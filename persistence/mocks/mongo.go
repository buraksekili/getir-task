package mocks

import (
	"context"
	"sync"

	"github.com/buraksekili/getir-task/persistence"
)

// MongoMock represents a mock structure for MongoDB.
type MongoMock struct {
	storage map[string]string
	mu      *sync.Mutex
}

// New returns a new mock db for MongoDB.
func New(storage map[string]string) persistence.Database {
	return &MongoMock{storage, &sync.Mutex{}}
}

func (im *MongoMock) FetchData(_ context.Context, f persistence.Filter) ([]persistence.FetchResObj, error) {
	im.mu.Lock()
	defer im.mu.Unlock()
	return []persistence.FetchResObj{{Key: f.Key, Val: im.storage[f.Key]}}, nil
}
