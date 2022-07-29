package stime

import "time"

// Get current utc unix timestamp in milliseconds
func UTCNowUnixMS() int64 {
	return time.Now().UTC().UnixMilli()
}

// UTC unix milliseconds timestamp to time
func MSToUTCTime(ms int64) time.Time {
	return time.UnixMilli(ms).UTC()
}
