package task

type ISliceScheduler interface {
	SliceRun(slicePtr interface{}, action func(i int, v interface{}))
	Cancel()
}
