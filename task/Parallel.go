package task

import (
	"reflect"

	log "github.com/kataras/golog"
	"github.com/syncfuture/go/dto"
)

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

	for i, action := range actions {
		go action(chs[i])
	}

	r = make([]*dto.ChannelResultDTO, actionCount)
	for i, ch := range chs {
		r[i] = <-ch
	}

	return
}

func (x *parallel) ForEach(slicePtr interface{}, action func(chan *dto.ChannelResultDTO, int, interface{})) (r []*dto.ChannelResultDTO) {
	v := reflect.ValueOf(slicePtr)
	if v.Kind() != reflect.Ptr {
		log.Fatal("slicePtr must be a slice pointer")
	}

	s := v.Elem()
	if s.Kind() != reflect.Slice {
		log.Fatal("slicePtr must be a slice pointer")
	}

	actionCount := s.Len()
	chs := make([]chan *dto.ChannelResultDTO, actionCount)
	for i := 0; i < actionCount; i++ {
		chs[i] = make(chan *dto.ChannelResultDTO)
	}

	for i := 0; i < s.Len(); i++ {
		go action(chs[i], i, s.Index(i).Interface())
	}

	r = make([]*dto.ChannelResultDTO, actionCount)
	for i, ch := range chs {
		r[i] = <-ch
	}

	return
}
