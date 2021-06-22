package sjson

import (
	"encoding/json"

	"github.com/syncfuture/go/serr"
	"github.com/syncfuture/go/u"
	"github.com/tidwall/gjson"
)

func UnmarshalSection(data []byte, path string, target interface{}) (err error) {
	v := gjson.GetBytes(data, path)
	err = json.Unmarshal(u.StrToBytes(v.Raw), target)
	return serr.WithStack(err)
}
