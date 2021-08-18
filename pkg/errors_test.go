package pkg

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	e1 := fmt.Errorf("error1")
	e2 := fmt.Errorf("test")
	e3 := Wrap(e1, e2)

	res := strings.Contains(e3.Error(), e1.Error())
	assert.Equal(t, true, res, fmt.Sprintf("expected true, got: %v", res))

	res = strings.Contains(e3.Error(), e2.Error())
	assert.Equal(t, true, res, fmt.Sprintf("expected true, got: %v", res))
}
