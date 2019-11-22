package task

import (
	"reflect"
	"sync"

	log "github.com/kataras/golog"
)

type flowScheduler struct {
	maxConcurrent int
}

// NewFlowScheduler create a flow scheduler
func NewFlowScheduler(maxConcurrent int) *flowScheduler {
	r := new(flowScheduler)
	r.maxConcurrent = maxConcurrent
	return r
}

// Run run with slice, make each action under scheduler control
func (x *flowScheduler) SliceRun(slicePtr interface{}, action func(i int, v interface{})) {
	v := reflect.ValueOf(slicePtr)
	if v.Kind() != reflect.Ptr {
		log.Fatal("slicePtr must be a slice pointer")
	}

	s := v.Elem()
	if s.Kind() != reflect.Slice {
		log.Fatal("slicePtr must be a slice pointer")
	}

	wg := sync.WaitGroup{}
	wg.Add(s.Len())
	ch := make(chan byte, x.maxConcurrent)

	for i := 0; i < s.Len(); i++ {
		ch <- 0
		go func(a int) {
			defer func() {
				<-ch
				wg.Done()
			}()
			action(a, s.Index(a).Interface())
		}(i)
	}

	wg.Wait()
}

// Run run under scheduler control
func (x *flowScheduler) Run(repeatTime int, action func(i int)) {
	wg := sync.WaitGroup{}
	wg.Add(repeatTime)
	ch := make(chan byte, x.maxConcurrent)

	for i := 0; i < repeatTime; i++ {
		ch <- 0
		go func(a int) {
			defer func() {
				<-ch
				wg.Done()
			}()
			action(a)
		}(i)
	}

	wg.Wait()
}
