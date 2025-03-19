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

// getFirstMondayOfYear 获取指定年份的第一个星期一
// 参数：
//   - year: 目标年份
//
// 返回：
//   - time.Time: 该年份第一个星期一的日期
func getFirstMondayOfYear(year int) time.Time {
	// 创建该年1月1日的日期对象
	firstDay := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	// 获取1月1日是星期几
	weekday := firstDay.Weekday()
	// 如果不是星期一，计算到第一个星期一的天数
	// 使用(time.Monday-weekday+7)%7确保结果为正数
	if weekday != time.Monday {
		firstDay = firstDay.AddDate(0, 0, int(time.Monday-weekday+7)%7)
	}
	return firstDay
}

// addWeeks 计算给定年份和周数加上或减去指定周数后的新年份和周数
// 参数：
//   - year: 起始年份
//   - week: 起始周数（1-53）
//   - deltaWeeks: 要增加或减少的周数（正数表示增加，负数表示减少）
//
// 返回：
//   - newYear: 计算后的新年份
//   - newWeek: 计算后的新周数（1-53）
func AddWeeks(year, week, deltaWeeks int) (newYear, newWeek int) {
	// 获取起始年份的第一周的第一个星期一作为基准日期
	yearFirstMonday := getFirstMondayOfYear(year)

	// 计算目标日期（targetWeekMonday）：
	// targetWeekMonday 表示从起始年份第一周开始，经过(week-1)周和deltaWeeks周后的具体日期
	// 由于firstDay是周一，且我们加上的天数都是7的倍数，所以targetWeekMonday也一定是周一
	// 例如：如果year=2024, week=3, deltaWeeks=2
	// 1. firstDay 是2024年第一个星期一（2024-01-01）
	// 2. (week-1)*7 计算从第一周到第三周的天数（2*7=14天）
	// 3. deltaWeeks*7 计算要增加的两周天数（2*7=14天）
	// 4. 最终 targetWeekMonday 就是 2024-01-01 + 28天 = 2024-01-29（周一）
	targetWeekMonday := yearFirstMonday.AddDate(0, 0, (week-1)*7+deltaWeeks*7)

	// 从目标日期中获取新的年份
	newYear = targetWeekMonday.Year()

	// 获取新年份的第一周的第一个星期一作为新的基准日期
	newYearFirstMonday := getFirstMondayOfYear(newYear)

	// 计算新的周数：
	// 1. targetWeekMonday.Sub(newYearFirstMonday): 计算目标日期与新年份第一周之间的时间差
	// 2. Hours()/24: 将时间差转换为天数
	// 3. days/7 + 1: 将天数转换为周数（加1是因为周数从1开始计数）
	days := int(targetWeekMonday.Sub(newYearFirstMonday).Hours() / 24)
	newWeek = (days / 7) + 1

	return
}

// 获取某一年的总周数（ISO 8601）
func GetWeeksInYear(year int) int {
	lastDay := time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC)
	_, week := lastDay.ISOWeek()
	return week
}
