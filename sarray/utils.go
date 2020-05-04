package sarray

func StrInArray(val string, array []string) bool {
	for _, item := range array {
		if item == val {
			return true
		}
	}
	return false
}

func IntInArray(val int, array []int) bool {
	for _, item := range array {
		if item == val {
			return true
		}
	}
	return false
}

func Int32InArray(val int32, array []int32) bool {
	for _, item := range array {
		if item == val {
			return true
		}
	}
	return false
}

func Int64InArray(val int64, array []int64) bool {
	for _, item := range array {
		if item == val {
			return true
		}
	}
	return false
}

func Float32InArray(val float32, array []float32) bool {
	for _, item := range array {
		if item == val {
			return true
		}
	}
	return false
}

func Float64InArray(val float64, array []float64) bool {
	for _, item := range array {
		if item == val {
			return true
		}
	}
	return false
}
