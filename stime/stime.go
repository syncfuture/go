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

// AddWeeks calculates the new year and week number after adding or subtracting a certain number of weeks
// to the given year and week number.
func AddWeeks(year, week, deltaWeeks int) (newYear, newWeek int) {
	// Create the first Monday of the current year
	firstDay := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	// Calculate the first Monday of the year
	for firstDay.Weekday() != time.Monday {
		firstDay = firstDay.AddDate(0, 0, 1)
	}

	// Calculate the start date of the given week
	startOfWeek := firstDay.AddDate(0, 0, (week-1)*7)

	// Add or subtract weeks
	newDate := startOfWeek.AddDate(0, 0, deltaWeeks*7)

	// Calculate the new year and week number
	newYear, newWeek = newDate.ISOWeek()

	return newYear, newWeek
}

// 获取某一年的总周数（ISO 8601）
func GetWeeksInYear(year int) int {
	lastDay := time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)
	_, week := lastDay.ISOWeek()
	return week
}
