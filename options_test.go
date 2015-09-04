package sol

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTag(t *testing.T) {
	tag, opts := parseTag(",omitempty,nullable")
	assert.Equal(t, "", tag)
	assert.Equal(t, options{OmitEmpty, "nullable"}, opts)
}
