package timeutil

import "time"

// 转换为字符串：yyyy-MM-dd HH:mm:ss
func ToString(t time.Time) string {
	return t.Format(time.DateTime)
}

// 转换为字符串：yyyy-MM-dd HH:mm:ss
func SecondToString(sec int64) string {
	t := GetTime(sec * TOTAL_MILLISECONDS_PER_SECOND)
	return t.Format(time.DateTime)
}

// 转换为字符串：yyyy-MM-dd HH:mm:ss
func MilliToString(ms int64) string {
	t := GetTime(ms)
	return t.Format(time.DateTime)
}

// 转换为字符串：yyyy-MM-dd
func ToDateString(t time.Time) string {
	return t.Format(time.DateOnly)
}

// 转换为字符串：yyyy-MM-dd
func SecondToDateString(sec int64) string {
	t := GetTime(sec * TOTAL_MILLISECONDS_PER_SECOND)
	return t.Format(time.DateOnly)
}

// 转换为字符串：yyyy-MM-dd
func MilliToDateString(ms int64) string {
	t := GetTime(ms)
	return t.Format(time.DateOnly)
}

// 转换为字符串：HH:mm:ss
func ToTimeString(t time.Time) string {
	return t.Format(time.TimeOnly)
}

// 转换为字符串：HH:mm:ss
func SecondToTimeString(sec int64) string {
	t := GetTime(sec * TOTAL_MILLISECONDS_PER_SECOND)
	return t.Format(time.TimeOnly)
}

// 转换为字符串：HH:mm:ss
func MilliToTimeString(ms int64) string {
	t := GetTime(ms)
	return t.Format(time.TimeOnly)
}
