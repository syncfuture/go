package sarray

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasAllStr(t *testing.T) {
	source := []string{"a", "b", "c", "d"}
	in1 := []string{"b", "e"}
	assert.False(t, HasAllStr(source, in1))

	in2 := []string{"b", "d"}
	assert.True(t, HasAllStr(source, in2))
}
