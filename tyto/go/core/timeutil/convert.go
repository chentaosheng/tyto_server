package timeutil

import "time"

// 将毫秒转换为time.Time
func GetTime(ms int64) time.Time {
	return time.Unix(0, ms*int64(time.Millisecond)).Local()
}

// 将time.Time转换为unix时间戳，单位: 秒
func GetSecond(t time.Time) int64 {
	return t.Unix()
}

// 将time.Time转换为unix时间戳，单位: 毫秒
func GetMilli(t time.Time) int64 {
	return t.UnixMilli()
}
