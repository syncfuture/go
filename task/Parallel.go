package task

import "github.com/syncfuture/go/dto"

type parallel struct{}

func NewParallel() *parallel {
	return new(parallel)
}

func (x *parallel) Invoke(actions []func(chan *dto.ChannelResultDTO)) (r []*dto.ChannelResultDTO) {
	actionCount := len(actions)
	chs := make([]chan *dto.ChannelResultDTO, actionCount)
	for i := 0; i < actionCount; i++ {
		chs[i] = make(chan *dto.ChannelResultDTO)
	}

	r = make([]*dto.ChannelResultDTO, actionCount)

	for i, action := range actions {
		go action(chs[i])
	}

	for i, ch := range chs {
		r[i] = <-ch
	}

	return
}
