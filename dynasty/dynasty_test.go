package dynasty

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartsWith(t *testing.T) {
	assert.True(t, startswith([]byte("bla"), []byte("blaireau")))
	assert.True(t, startswith([]byte("bla"), []byte("bla")))
	assert.False(t, startswith([]byte("blaireau"), []byte("bla")))
}
