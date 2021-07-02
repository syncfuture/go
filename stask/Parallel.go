package stask

import (
	"reflect"

	log "github.com/kataras/golog"
	"github.com/syncfuture/go/sdto"
)

type parallel struct{}

func NewParallel() *parallel {
	return new(parallel)
}

func (x *parallel) Invoke(actions ...func(chan *sdto.ChannelResultDTO)) (r []*sdto.ChannelResultDTO) {
	actionCount := len(actions)
	chs := make([]chan *sdto.ChannelResultDTO, actionCount)
	for i := 0; i < actionCount; i++ {
		chs[i] = make(chan *sdto.ChannelResultDTO)
	}

	for i, action := range actions {
		go action(chs[i])
	}

	r = make([]*sdto.ChannelResultDTO, actionCount)
	for i, ch := range chs {
		r[i] = <-ch
	}

	return
}

func (x *parallel) ForEach(slicePtr interface{}, action func(chan *sdto.ChannelResultDTO, int, interface{})) (r []*sdto.ChannelResultDTO) {
	v := reflect.ValueOf(slicePtr)
	if v.Kind() != reflect.Ptr {
		log.Fatal("slicePtr must be a slice pointer")
	}

	s := v.Elem()
	if s.Kind() != reflect.Slice {
		log.Fatal("slicePtr must be a slice pointer")
	}

	actionCount := s.Len()
	chs := make([]chan *sdto.ChannelResultDTO, actionCount)
	for i := 0; i < actionCount; i++ {
		chs[i] = make(chan *sdto.ChannelResultDTO)
	}

	for i := 0; i < s.Len(); i++ {
		go action(chs[i], i, s.Index(i).Interface())
	}

	r = make([]*sdto.ChannelResultDTO, actionCount)
	for i, ch := range chs {
		r[i] = <-ch
	}

	return
}
