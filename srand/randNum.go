package srand

import (
	"math/rand"
	"time"

	log "github.com/syncfuture/go/slog"
)

func IntRange(min, max int) int {
	if min > max {
		log.Fatal("min cannot greater than max")
	} else if min == max {
		return max
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func Int31Range(min, max int32) int32 {
	if min > max {
		log.Fatal("min cannot greater than max")
	} else if min == max {
		return max
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Int31n(max-min) + min
}

func Int63Range(min, max int64) int64 {
	if min > max {
		log.Fatal("min cannot greater than max")
	} else if min == max {
		return max
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min) + min
}
