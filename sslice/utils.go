package sslice

import "math"

// Generic type constraint
type comparable interface {
	string | int | int32 | int64
}

// Private generic function to check if an element exists in a slice
func has[T comparable](source []T, in T) bool {
	sourceMap := make(map[T]struct{}, len(source))
	for _, item := range source {
		sourceMap[item] = struct{}{}
	}
	_, found := sourceMap[in]
	return found
}

// Private generic function to check if any element from a slice exists in another slice
func hasAny[T comparable](source []T, in []T) bool {
	sourceMap := make(map[T]struct{}, len(source))
	for _, item := range source {
		sourceMap[item] = struct{}{}
	}
	for _, v := range in {
		if _, found := sourceMap[v]; found {
			return true
		}
	}
	return false
}

// Private generic function to check if all elements from a slice exist in another slice
func hasAll[T comparable](source []T, in []T) bool {
	sourceMap := make(map[T]struct{}, len(source))
	for _, item := range source {
		sourceMap[item] = struct{}{}
	}
	for _, v := range in {
		if _, found := sourceMap[v]; !found {
			return false
		}
	}
	return true
}

// Function to handle floating-point number precision issues
func hasFloat[T float64 | float32](source []T, in T, epsilon T) bool {
	for _, item := range source {
		if math.Abs(float64(item-in)) < float64(epsilon) {
			return true
		}
	}
	return false
}

// Function to handle hasAny for float types
func hasAnyFloat[T float64 | float32](source []T, in []T, epsilon T) bool {
	for _, v := range in {
		if hasFloat(source, v, epsilon) {
			return true
		}
	}
	return false
}

// Function to handle hasAll for float types
func hasAllFloat[T float64 | float32](source []T, in []T, epsilon T) bool {
	for _, v := range in {
		if !hasFloat(source, v, epsilon) {
			return false
		}
	}
	return true
}

// Public functions
func HasStr(source []string, in string) bool      { return has(source, in) }
func HasAnyStr(source []string, in []string) bool { return hasAny(source, in) }
func HasAllStr(source []string, in []string) bool { return hasAll(source, in) }

func HasInt(source []int, in int) bool      { return has(source, in) }
func HasAnyInt(source []int, in []int) bool { return hasAny(source, in) }
func HasAllInt(source []int, in []int) bool { return hasAll(source, in) }

func HasInt32(source []int32, in int32) bool      { return has(source, in) }
func HasAnyInt32(source []int32, in []int32) bool { return hasAny(source, in) }
func HasAllInt32(source []int32, in []int32) bool { return hasAll(source, in) }

func HasInt64(source []int64, in int64) bool      { return has(source, in) }
func HasAnyInt64(source []int64, in []int64) bool { return hasAny(source, in) }
func HasAllInt64(source []int64, in []int64) bool { return hasAll(source, in) }

func HasFloat64(source []float64, in float64) bool      { return hasFloat(source, in, 1e-9) }
func HasAnyFloat64(source []float64, in []float64) bool { return hasAnyFloat(source, in, 1e-9) }
func HasAllFloat64(source []float64, in []float64) bool { return hasAllFloat(source, in, 1e-9) }

func HasFloat32(source []float32, in float32) bool      { return hasFloat(source, in, 1e-6) }
func HasAnyFloat32(source []float32, in []float32) bool { return hasAnyFloat(source, in, 1e-6) }
func HasAllFloat32(source []float32, in []float32) bool { return hasAllFloat(source, in, 1e-6) }
