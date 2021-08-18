package inmemory

import (
	"fmt"
	"testing"

	"github.com/buraksekili/getir-task/persistence"
	"github.com/stretchr/testify/assert"
)

func TestFetchData(t *testing.T) {
	im := New(map[string]string{"key": "value"})
	res := im.FetchData("key")
	assert.Equal(t, "value", res, fmt.Sprintf("expected 'value', got: %s", res))

	res = im.FetchData("keyinv")
	assert.Equal(t, "", res, fmt.Sprintf("expected '', got: %s", res))
}

func TestStoreData(t *testing.T) {
	im := New(map[string]string{})
	im.StoreData(persistence.InMemoryKeyValue{Key: "key", Value: "value"})

	res := im.FetchData("key")
	assert.Equal(t, "value", res, fmt.Sprintf("expected 'value', got: %s", res))
}
