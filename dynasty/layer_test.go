package dynasty

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNum(t *testing.T) {
	assert.Equal(t, []int{0, 1}, num(26, []int{}))
	assert.Equal(t, []int{1, 2}, num(53, []int{}))
}

func TestEncode(t *testing.T) {
	assert.Equal(t, []byte("a__"), encode(0, 3))
	assert.Equal(t, []byte("z__"), encode(25, 3))
	assert.Equal(t, []byte("ab_"), encode(26, 3))
	assert.Equal(t, []byte("bc_"), encode(53, 3))
}
