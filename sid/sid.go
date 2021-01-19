package sid

import (
	"fmt"

	"github.com/sony/sonyflake"
	log "github.com/syncfuture/go/slog"
)

type IIDGenerator interface {
	GenerateString() string
}

type SonyflakeIDGenerator struct {
	generator *sonyflake.Sonyflake
}

func NewSonyflakeIDGenerator() IIDGenerator {
	return &SonyflakeIDGenerator{
		generator: sonyflake.NewSonyflake(sonyflake.Settings{}),
	}
}

func (x *SonyflakeIDGenerator) GenerateString() string {
	id, err := x.generator.NextID()
	if err != nil {
		log.Errorf("flake.NextID() failed with %s\n", err)
	}
	return fmt.Sprintf("%x", id)
}
