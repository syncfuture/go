package config

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	_json = []byte(`{
		"Dev":{
			"Debug":true
		},
		"Count": 79,
		"L1":{
			"L2":{
				"L3":{
					"L4":"ABC"
				}
			}
		}
	}`)
	_config *Configuration
)

func init() {
	e := make(Configuration)

	// Unmarshal json data structure
	err := json.Unmarshal(_json, &e)
	if err != nil {
		panic(err)
	}

	_config = &e
}

func TestGetString(t *testing.T) {
	r := _config.GetString("L1.L2.L3.L4")
	assert.Equal(t, "ABC", r)
}

func TestGetBool(t *testing.T) {
	r := _config.GetBool("Dev.Debug")
	assert.Equal(t, true, r)
}

func TestGetInt(t *testing.T) {
	r := _config.GetInt("Count")
	assert.Equal(t, 79, r)
}
