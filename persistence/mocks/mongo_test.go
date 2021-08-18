package mocks

import (
	"context"
	"fmt"
	"testing"

	"github.com/buraksekili/getir-task/persistence"
	"github.com/stretchr/testify/assert"
)

func TestFetchData(t *testing.T) {
	db := New(map[string]string{"key": "value"})
	res, err := db.FetchData(context.Background(), persistence.Filter{Key: "key"})
	assert.Nil(t, err, fmt.Sprintf("unexpected error during fetching data %s", err))
	assert.Equal(t, 1, len(res), fmt.Sprintf("expected length=1, got %d", len(res)))
	assert.Equal(t, "value", res[0].Val, fmt.Sprintf("expected 'value', got %s", res[0].Val))
}
