package handlers

import (
	"github.com/buraksekili/getir-task/persistence"
)

// FetchDataRes represents a result JSON for fetch data requests.
type FetchDataRes struct {
	Code    int                       `json:"code"`
	Msg     string                    `json:"msg"`
	Records []persistence.FetchResObj `json:"records"`
}

// ErrorRes represents a JSON response for failures.
type ErrorRes struct {
	Error string `json:"error"`
}

// InMemoryRes represents a JSON response for in-memory requests.
type InMemoryRes struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
