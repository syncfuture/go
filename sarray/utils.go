package sarray

func HasStr(array []string, in string) bool {
	for _, item := range array {
		if item == in {
			return true
		}
	}
	return false
}

func HasAnyStr(array []string, in []string) bool {
	for _, item := range array {
		for _, v := range in {
			if item == v {
				return true
			}
		}
	}
	return false
}

func HasInt(array []int, in int) bool {
	for _, item := range array {
		if item == in {
			return true
		}
	}
	return false
}

func HasAnyInt(array []int, in []int) bool {
	for _, item := range array {
		for _, v := range in {
			if item == v {
				return true
			}
		}
	}
	return false
}

func HasInt32(array []int32, in int32) bool {
	for _, item := range array {
		if item == in {
			return true
		}
	}
	return false
}

func HasAnyInt32(array []int32, in []int32) bool {
	for _, item := range array {
		for _, v := range in {
			if item == v {
				return true
			}
		}
	}
	return false
}

func HasInt64(array []int64, in int64) bool {
	for _, item := range array {
		if item == in {
			return true
		}
	}
	return false
}

func HasAnyInt64(array []int64, in []int64) bool {
	for _, item := range array {
		for _, v := range in {
			if item == v {
				return true
			}
		}
	}
	return false
}

func HasFloat32(array []float32, in float32) bool {
	for _, item := range array {
		if item == in {
			return true
		}
	}
	return false
}

func HasAnyFloat32(array []float32, in []float32) bool {
	for _, item := range array {
		for _, v := range in {
			if item == v {
				return true
			}
		}
	}
	return false
}

func HasFloat64(array []float64, in float64) bool {
	for _, item := range array {
		if item == in {
			return true
		}
	}
	return false
}

func HasAnyFloat64(array []float64, in []float64) bool {
	for _, item := range array {
		for _, v := range in {
			if item == v {
				return true
			}
		}
	}
	return false
}
