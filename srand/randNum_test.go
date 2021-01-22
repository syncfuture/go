package srand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntRange(t *testing.T) {
	min := -50
	max := -2
	a := IntRange(min, max)
	b := IntRange(min, max)

	t.Log(a, b)
	assert.NotEqual(t, a, b)
}
