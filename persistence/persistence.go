package persistence

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrAggregation shows failure during aggregation.
	ErrAggregation = errors.New("failed to aggregate")
)

// Filter represents required structure to filter MongoDB queries.
type Filter struct {
	StartDate     string `json:"startDate" bson:"startDate"`
	EndDate       string `json:"endDate" bson:"endDate"`
	MaxTotalCount int    `json:"maxCount" bson:"minCount"`
	MinTotalCount int    `json:"minCount" bson:"maxCount"`
	Key           string `json:"key,omitempty"`
}

// FetchResObj represents a single document in MongoDB Collection.
type FetchResObj struct {
	Key        string    `json:"key" bson:"key"`
	Val        string    `json:"val,omitempty"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
	Counts     []int64   `json:"counts,omitempty" bson:"counts"`
	TotalCount int       `json:"totalCount" bson:"totalCount"`
}

// InMemoryKeyValue represents a key-value pair for in-memory database.
type InMemoryKeyValue struct {
	Key   string
	Value string
}

// Database represents an interface for Databases.
type Database interface {
	FetchData(ctx context.Context, f Filter) ([]FetchResObj, error)
}
