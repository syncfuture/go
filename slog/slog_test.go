package slog

import (
	"os"
	"testing"
	"time"

	"github.com/syncfuture/go/sconfig"
	"github.com/syncfuture/go/serr"
)

func TestDebug(t *testing.T) {
	cp := sconfig.NewJsonConfigProvider()
	Init(cp)

	for i := 0; i < 100; i++ {
		Debug(time.Now())
		time.Sleep(1 * time.Second)
	}
}

func TestErr(t *testing.T) {
	err := test1()
	Warn(err)
}

func test1() error {
	return test2()
}

func test2() error {
	return test3()
}

func test3() error {
	_, err := os.Open("test.aaa")
	return serr.Wrap(err)
}
