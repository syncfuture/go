package scheduler

import (
	"reflect"
	"sync"

	log "github.com/kataras/golog"
)

type goroutinesScheduler struct {
	maxConcurrent int
}

// NewGoroutineScheduler create a gorutine scheduler
func NewGoroutineScheduler(maxConcurrent int) *goroutinesScheduler {
	r := new(goroutinesScheduler)
	r.maxConcurrent = maxConcurrent
	return r
}

// Run batch run slice, process each item under scheduler control
func (x *goroutinesScheduler) RunSlice(slicePtr interface{}, proccessAction func(*sync.WaitGroup, chan byte, reflect.Value)) {
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
		go proccessAction(&wg, ch, s.Index(i))
	}
}

// Run batch run, process under scheduler control
func (x *goroutinesScheduler) Run(repeatTime int, proccessAction func(*sync.WaitGroup, chan byte, int)) {
	wg := sync.WaitGroup{}
	wg.Add(repeatTime)
	ch := make(chan byte, x.maxConcurrent)

	for i := 0; i < repeatTime; i++ {
		ch <- 0
		go proccessAction(&wg, ch, i)
	}
}
