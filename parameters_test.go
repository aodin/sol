package sol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParameters_Add(t *testing.T) {
	ps := Params()
	ps.Add(1)
	assert.Equal(t, 1, ps.Len())
}
