package slog

import (
	"testing"
	"time"

	"github.com/syncfuture/go/config"
)

func Test1(t *testing.T) {
	cp := config.NewJsonConfigProvider()
	Init(cp)

	for {
		Debug(time.Now())
		time.Sleep(1 * time.Second)
	}
}
