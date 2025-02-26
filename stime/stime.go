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

func GetWeekStartEnd(year, week int) (time.Time, time.Time) {
	// Get the first day of the year
	firstDay := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	// Calculate the weekday of the first day
	weekday := int(firstDay.Weekday())
	if weekday == 0 {
		weekday = 7 // Set Sunday to 7
	}

	// Calculate the first day of the given week (ISO weeks start on Monday)
	startDate := firstDay.AddDate(0, 0, (week-1)*7-(weekday-1))
	endDate := startDate.AddDate(0, 0, 6) // Calculate Sunday

	return startDate, endDate
}
