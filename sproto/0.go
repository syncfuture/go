package sproto

import (
	"encoding/json"

	"github.com/syncfuture/go/serr"
)

func GetResultData[T any](in *Result) (T, error) {
	var r T

	err := json.Unmarshal(in.Bytes, &r)

	if err != nil {
		return r, serr.WithStack(err)
	}

	return r, nil
}

func NewResult[T any](in T) (r *Result, err error) {
	r = new(Result)

	r.Bytes, err = json.Marshal(in)
	if err != nil {
		return nil, serr.WithStack(err)
	}

	return r, nil
}
