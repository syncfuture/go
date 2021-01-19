package slog

import (
	"testing"
	"time"

	"github.com/syncfuture/go/sconfig"
)

func Test1(t *testing.T) {
	cp := sconfig.NewJsonConfigProvider()
	Init(cp)

	for i := 0; i < 100; i++ {
		Debug(time.Now())
		time.Sleep(1 * time.Second)
	}
}
