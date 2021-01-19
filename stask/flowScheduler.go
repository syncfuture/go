package stask

import (
	"reflect"
	"sync"

	log "github.com/syncfuture/go/slog"
)

type flowScheduler struct {
	maxConcurrent int
	cancel        bool
	rwMutex       *sync.RWMutex
}

// NewFlowScheduler create a flow scheduler
func NewFlowScheduler(maxConcurrent int) *flowScheduler {
	r := new(flowScheduler)
	r.maxConcurrent = maxConcurrent
	r.rwMutex = new(sync.RWMutex)
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

		x.rwMutex.RLock()
		if x.cancel {
			x.rwMutex.RUnlock()
			return
		}
		x.rwMutex.RUnlock()

		go func(a int) {
			defer func() {
				<-ch
				wg.Done()
			}()

			x.rwMutex.RLock()
			if x.cancel {
				x.rwMutex.RUnlock()
				return
			}
			x.rwMutex.RUnlock()

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

func (x *flowScheduler) Cancel() {
	x.rwMutex.Lock()
	x.cancel = true
	x.rwMutex.Unlock()
}
