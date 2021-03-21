package stask

import (
	"testing"

	"github.com/syncfuture/go/sdto"
	"github.com/syncfuture/go/serr"
	"github.com/syncfuture/go/slog"
)

func Test(t *testing.T) {
	p := NewParallel()
	rs := p.Invoke([]func(chan *sdto.ChannelResultDTO){
		func(ch chan *sdto.ChannelResultDTO) { test1(ch, 1) },
		func(ch chan *sdto.ChannelResultDTO) { test2(ch, "abc") },
		func(ch chan *sdto.ChannelResultDTO) { test3(ch, false) },
	})

	for _, r := range rs {
		slog.Info(r.Result)
	}
}

func test1(ch chan *sdto.ChannelResultDTO, arg int) {
	slog.Info("test1 running...")

	defer func() {
		ch <- &sdto.ChannelResultDTO{
			Result: arg,
			Error:  serr.New("test1 error"),
		}
	}()
}

func test2(ch chan *sdto.ChannelResultDTO, arg string) {
	slog.Info("test2 running...")

	defer func() {
		ch <- &sdto.ChannelResultDTO{
			Result: arg,
			Error:  serr.New("test2 error"),
		}
	}()
}

func test3(ch chan *sdto.ChannelResultDTO, arg bool) {
	slog.Info("test3 running...")

	defer func() {
		ch <- &sdto.ChannelResultDTO{
			Result: arg,
			Error:  serr.New("test3 error"),
		}
	}()
}
