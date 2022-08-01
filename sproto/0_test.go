package sproto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractResultData(t *testing.T) {
	address := &Address{
		Address1: "123 Aaa Road",
		City:     "Bbb",
		State:    "CA",
		ZipCode:  "22991",
		Country:  "US",
	}

	rs, err := NewResult(address)
	assert.NoError(t, err)
	assert.NotNil(t, rs)
	assert.NotEmpty(t, rs.Bytes)
}
