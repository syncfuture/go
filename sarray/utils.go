package sarray

func HasStr(source []string, in string) bool {
	for _, item := range source {
		if item == in {
			return true
		}
	}
	return false
}
func HasAnyStr(source []string, in []string) bool {
	for _, v := range in {
		for _, item := range source {
			if item == v {
				return true
			}
		}
	}
	return false
}
func HasAllStr(source []string, in []string) bool {
	for _, v := range in {
		var found bool
		for _, item := range source {
			if item == v {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}
	return true
}

func HasInt(source []int, in int) bool {
	for _, item := range source {
		if item == in {
			return true
		}
	}
	return false
}
func HasAnyInt(source []int, in []int) bool {
	for _, v := range in {
		for _, item := range source {
			if item == v {
				return true
			}
		}
	}
	return false
}
func HasAllInt(source []int, in []int) bool {
	for _, v := range in {
		var found bool
		for _, item := range source {
			if item == v {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}
	return true
}

func HasInt32(source []int32, in int32) bool {
	for _, item := range source {
		if item == in {
			return true
		}
	}
	return false
}
func HasAnyInt32(source []int32, in []int32) bool {
	for _, v := range in {
		for _, item := range source {
			if item == v {
				return true
			}
		}
	}
	return false
}
func HasAllInt32(source []int32, in []int32) bool {
	for _, v := range in {
		var found bool
		for _, item := range source {
			if item == v {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}
	return true
}

func HasInt64(source []int64, in int64) bool {
	for _, item := range source {
		if item == in {
			return true
		}
	}
	return false
}
func HasAnyInt64(source []int64, in []int64) bool {
	for _, v := range in {
		for _, item := range source {
			if item == v {
				return true
			}
		}
	}
	return false
}
func HasAllInt64(source []int64, in []int64) bool {
	for _, v := range in {
		var found bool
		for _, item := range source {
			if item == v {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}
	return true
}

func HasFloat32(source []float32, in float32) bool {
	for _, item := range source {
		if item == in {
			return true
		}
	}
	return false
}
func HasAnyFloat32(source []float32, in []float32) bool {
	for _, v := range in {
		for _, item := range source {
			if item == v {
				return true
			}
		}
	}
	return false
}
func HasAllFloat32(source []float32, in []float32) bool {
	for _, v := range in {
		var found bool
		for _, item := range source {
			if item == v {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}
	return true
}

func HasFloat64(source []float64, in float64) bool {
	for _, item := range source {
		if item == in {
			return true
		}
	}
	return false
}
func HasAnyFloat64(source []float64, in []float64) bool {
	for _, v := range in {
		for _, item := range source {
			if item == v {
				return true
			}
		}
	}
	return false
}
func HasAllFloat64(source []float64, in []float64) bool {
	for _, v := range in {
		var found bool
		for _, item := range source {
			if item == v {
				found = true
				break
			}
		}

		if !found {
			return false
		}
	}
	return true
}
