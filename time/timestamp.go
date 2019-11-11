package time

import (
	t "time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	log "github.com/kataras/golog"
)

func TimestampUTCNow() *timestamp.Timestamp {
	ts, err := ptypes.TimestampProto(t.Now().UTC())
	if err != nil {
		log.Error("ptypes: time.Now().UTC() out of Timestamp range")
	}
	return ts
}
