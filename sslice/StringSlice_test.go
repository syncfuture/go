package sslice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringSlice(t *testing.T) {
	x := make(StringSlice, 5)
	x[0] = "a"
	x[1] = "b"
	x[2] = "c"
	x[3] = "d"
	x[4] = "e"

	assert.True(t, x.Has("c"))

	x.Add("f")
	assert.True(t, x.Has("f"))

	x.Remove("c")
	x.Remove("e")
	assert.Len(t, x, 4)

	assert.False(t, x.Has("c"))
}
