package task

import (
	"errors"
	"testing"

	"github.com/syncfuture/go/dto"
	"github.com/syncfuture/go/slog"
)

func Test(t *testing.T) {
	p := NewParallel()
	rs := p.Invoke([]func(chan *dto.ChannelResultDTO){
		func(ch chan *dto.ChannelResultDTO) { test1(ch, 1) },
		func(ch chan *dto.ChannelResultDTO) { test2(ch, "abc") },
		func(ch chan *dto.ChannelResultDTO) { test3(ch, false) },
	})

	for _, r := range rs {
		slog.Info(r.Result)
	}
}

func test1(ch chan *dto.ChannelResultDTO, arg int) {
	slog.Info("test1 running...")

	defer func() {
		ch <- &dto.ChannelResultDTO{
			Result: arg,
			Error:  errors.New("test1 error"),
		}
	}()
}

func test2(ch chan *dto.ChannelResultDTO, arg string) {
	slog.Info("test2 running...")

	defer func() {
		ch <- &dto.ChannelResultDTO{
			Result: arg,
			Error:  errors.New("test2 error"),
		}
	}()
}

func test3(ch chan *dto.ChannelResultDTO, arg bool) {
	slog.Info("test3 running...")

	defer func() {
		ch <- &dto.ChannelResultDTO{
			Result: arg,
			Error:  errors.New("test3 error"),
		}
	}()
}
