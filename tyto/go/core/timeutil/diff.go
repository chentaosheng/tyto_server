package timeutil

// 判断是否为同一天，参数为毫秒时间戳
func IsSameDay(ms1, ms2 int64) bool {
	t1 := GetTime(ms1)
	t2 := GetTime(ms2)
	return t1.Year() == t2.Year() && t1.YearDay() == t2.YearDay()
}

// 判断是否为同一天，按服务器每日重置时间判断，参数为毫秒时间戳
func IsSameDailyDay(ms1 int64, ms2 int64, resetHour int64, resetMinute int64) bool {
	offset := resetHour*TOTAL_MILLISECONDS_PER_HOUR + resetMinute*TOTAL_MILLISECONDS_PER_MINUTE
	t1 := GetTime(ms1 - offset)
	t2 := GetTime(ms2 - offset)
	return t1.Year() == t2.Year() && t1.YearDay() == t2.YearDay()
}

// 判断是否为昨天，参数为毫秒时间戳
// yesterday 必须小于 today
func IsYesterday(yesterday, today int64) bool {
	t1 := GetTime(yesterday)
	t2 := GetTime(today)
	return t1.Year() == t2.Year() && t1.YearDay() == t2.YearDay()-1
}

// 判断是否为同一周，参数为毫秒时间戳
func IsSameWeek(ms1, ms2 int64) bool {
	t1 := GetTime(ms1)
	t2 := GetTime(ms2)
	_, w1 := t1.ISOWeek()
	_, w2 := t2.ISOWeek()
	return t1.Year() == t2.Year() && w1 == w2
}

// 判断是否为同一月，参数为毫秒时间戳
func IsSameMonth(ms1, ms2 int64) bool {
	t1 := GetTime(ms1)
	t2 := GetTime(ms2)
	return t1.Year() == t2.Year() && t1.Month() == t2.Month()
}

// ms2和ms1之间相差的天数
// ms2必须大于ms1，否则返回负数
func DaysBetween(ms1, ms2 int64) int32 {
	t1 := BeginningOfDayMilli(ms1)
	t2 := BeginningOfDayMilli(ms2)
	return int32((t2 - t1) / TOTAL_MILLISECONDS_PER_DAY)
}

// ms2和ms1之间相差的天数，按服务器每日重置时间判断
// ms2必须大于ms1，否则返回负数
func DailyDaysBetween(ms1, ms2 int64, resetHour int64, resetMinute int64) int32 {
	offset := resetHour*TOTAL_MILLISECONDS_PER_HOUR + resetMinute*TOTAL_MILLISECONDS_PER_MINUTE
	t1 := BeginningOfDayMilli(ms1 - offset)
	t2 := BeginningOfDayMilli(ms2 - offset)
	return int32((t2 - t1) / TOTAL_MILLISECONDS_PER_DAY)
}
