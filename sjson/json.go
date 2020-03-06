package sjson

import (
	"encoding/json"

	"github.com/tidwall/gjson"
)

func UnmarshalSection(data []byte, path string, target interface{}) (err error) {
	v := gjson.GetBytes(data, path)
	err = json.Unmarshal([]byte(v.Raw), target)
	return err
}
