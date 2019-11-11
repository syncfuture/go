package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type configDTO struct {
	ABC int
}

func TestJsonConfigProvider_GetConfig(t *testing.T) {
	a := configDTO{}
	b := JsonConfigProvider{}
	b.GetConfig(&a, "configs.json")
	assert.Equal(t, 123, a.ABC)
}
