package siris

import (
	"time"

	"github.com/syncfuture/go/sproto/timestamp"
)

func TimeStamp(t *timestamp.Timestamp, args ...interface{}) string {
	utcTime, _ := t.Time()
	beijingTime := utcTime.Add(8 * time.Hour)
	if len(args) > 0 {
		return beijingTime.Format(args[0].(string))
	} else {
		return beijingTime.Format("01/02/2006 03:04:05 PM")
	}
}

func Iter(count int32) []int32 {
	var i int32
	var r []int32
	for i = 0; i < count; i++ {
		r = append(r, i)
	}
	return r
}
