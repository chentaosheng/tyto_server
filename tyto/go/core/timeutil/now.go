package timeutil

import "time"

// 获取当前时间，时区为系统时区
func Now() time.Time {
	return time.Now().Local()
}

// 获取unix时间戳，单位秒
// unix时间戳是指从1970年1月1日（UTC/GMT的午夜）到现在的秒数
func NowSecond() int64 {
	return time.Now().Unix()
}

// 获取unix时间戳，单位毫秒
func NowMilli() int64 {
	return time.Now().UnixMilli()
}
