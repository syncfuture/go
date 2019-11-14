package stime

import (
	t "time"

	"github.com/syncfuture/go/sproto/timestamp"

	log "github.com/kataras/golog"
)

func TimestampUTCNow() *timestamp.Timestamp {
	ts, err := TimestampProto(t.Now().UTC())
	if err != nil {
		log.Error("ptypes: time.Now().UTC() out of Timestamp range")
	}
	return ts
}
