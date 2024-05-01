package timeutil

import "time"

// 获取小时的开始时间
func BeginningOfHour(t time.Time) time.Time {
	y, m, d := t.Date()
	h := t.Hour()
	return time.Date(y, m, d, h, 0, 0, 0, t.Location())
}

func BeginningOfHourMilli(ms int64) int64 {
	t := GetTime(ms)
	begin := BeginningOfHour(t)
	return GetMilli(begin)
}

func BeginningOfHourSecond(sec int64) int64 {
	return BeginningOfHourMilli(sec*TOTAL_MILLISECONDS_PER_SECOND) / TOTAL_MILLISECONDS_PER_SECOND
}

// 获取t的当天开始时间，即凌晨
func BeginningOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func BeginningOfDayMilli(ms int64) int64 {
	t := GetTime(ms)
	begin := BeginningOfDay(t)
	return GetMilli(begin)
}

func BeginningOfDaySecond(sec int64) int64 {
	return BeginningOfDayMilli(sec*TOTAL_MILLISECONDS_PER_SECOND) / TOTAL_MILLISECONDS_PER_SECOND
}

// 获取t的当周开始时间，从周一开始算
func BeginningOfWeek(t time.Time) time.Time {
	y, m, d := t.Date()
	wd := int(t.Weekday())
	if wd == 0 {
		wd = 7
	}
	return time.Date(y, m, d-wd+1, 0, 0, 0, 0, t.Location())
}

func BeginningOfWeekMilli(ms int64) int64 {
	t := GetTime(ms)
	begin := BeginningOfWeek(t)
	return GetMilli(begin)
}

func BeginningOfWeekSecond(sec int64) int64 {
	return BeginningOfWeekMilli(sec*TOTAL_MILLISECONDS_PER_SECOND) / TOTAL_MILLISECONDS_PER_SECOND
}

// 获取t的当月开始时间
func BeginningOfMonth(t time.Time) time.Time {
	y, m, _ := t.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, t.Location())
}

func BeginningOfMonthMilli(ms int64) int64 {
	t := GetTime(ms)
	begin := BeginningOfMonth(t)
	return GetMilli(begin)
}

func BeginningOfMonthSecond(sec int64) int64 {
	return BeginningOfMonthMilli(sec*TOTAL_MILLISECONDS_PER_SECOND) / TOTAL_MILLISECONDS_PER_SECOND
}
