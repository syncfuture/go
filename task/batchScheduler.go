package task

import (
	"math"
	"reflect"
	"sync"
	"time"

	log "github.com/syncfuture/go/slog"
)

type batchScheduler struct {
	batchSize int
	// Intervel intervel milliseconds per batch
	intervalMS  int
	cancel      bool
	onBatchDone func(int, int)
}

func NewBatchScheduler(batchSize, intervalMS int, batchEvents ...func(int, int)) *batchScheduler {
	r := &batchScheduler{
		batchSize:  batchSize,
		intervalMS: intervalMS,
	}
	if len(batchEvents) > 0 {
		r.onBatchDone = batchEvents[0]
	}
	return r
}
func (x *batchScheduler) SliceRun(slicePtr interface{}, action func(i int, v interface{})) {
	v := reflect.ValueOf(slicePtr)
	if v.Kind() != reflect.Ptr {
		log.Fatal("slicePtr must be a slice pointer")
	}

	s := v.Elem()
	if s.Kind() != reflect.Slice {
		log.Fatal("slicePtr must be a slice pointer")
	}

	if x.batchSize <= 0 {
		log.Fatal("batch size must be positive number.")
	}

	if x.intervalMS < 0 {
		log.Fatal("interval must be 0 or positive number.")
	}

	totalCount := s.Len()
	// 总页数
	totalPages := int(math.Ceil(float64(totalCount) / float64(x.batchSize)))

	wg := &sync.WaitGroup{}
	for pageIndex := 0; pageIndex < totalPages; pageIndex++ {
		if pageIndex == totalPages-1 {
			remainder := totalCount % x.batchSize
			if remainder == 0 {
				wg.Add(x.batchSize) // 页数
			} else {
				wg.Add(remainder) // 余下的个数
			}
		} else {
			wg.Add(x.batchSize) // 页数
		}

		for pageLoopIndex := 0; pageLoopIndex < x.batchSize; pageLoopIndex++ {
			itemIndex := pageIndex*x.batchSize + pageLoopIndex
			if itemIndex >= totalCount {
				break
			}
			go func(i int) {
				defer wg.Done()
				action(i, s.Index(i).Interface())
			}(itemIndex)
		}

		wg.Wait()
		if x.cancel {
			break
		}

		if x.onBatchDone != nil {
			// 触发批次执行完毕事件
			x.onBatchDone(pageIndex, totalPages)
		}
		if x.intervalMS > 0 && pageIndex < totalPages-1 {
			time.Sleep(time.Duration(x.intervalMS) * time.Millisecond)
		}
	}
}

func (x *batchScheduler) Cancel() {
	x.cancel = true
}
