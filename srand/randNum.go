package srand

import "math/rand"

func IntRange(min, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Intn(max-min) + min
}

func Int31Range(min, max int32) int32 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int31n(max-min) + min
}

func Int63Range(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}
