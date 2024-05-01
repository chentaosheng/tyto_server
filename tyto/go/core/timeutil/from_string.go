package timeutil

import (
	"strconv"
	"strings"
	"time"
	"tyto/core/tyto"
)

// 将yyyy-MM-dd HH:mm:ss格式的字符串转换为time.Time
func StringToTime(ctx tyto.Context, str string) (time.Time, bool) {
	t, err := time.ParseInLocation(time.DateTime, str, time.Local)
	if err != nil {
		ctx.Logger().Error("convert string to time failed:", err)
		return t, false
	}

	return t, true
}

// 将yyyy-MM-dd HH:mm:ss格式的字符串转换为1970开始的毫秒数
func StringToMilli(ctx tyto.Context, str string) (int64, bool) {
	t, ok := StringToTime(ctx, str)
	if !ok {
		return 0, false
	}
	return GetMilli(t), true
}

// 将yyyy-MM-dd HH:mm:ss格式的字符串转换为1970开始的秒数
func StringToSecond(ctx tyto.Context, str string) (int64, bool) {
	t, ok := StringToTime(ctx, str)
	if !ok {
		return 0, false
	}
	return GetSecond(t), true
}

// 将1M 1w 12d 23h 55m 00s格式的字符串转换为秒数
// 当带有M时，需要传入baseTime，通过baseTime所在月份计算出实际的天数
// 格式说明：
//  1. M = month, w = week, d = day, h = hour, m = minute, s = second
//  2. 数值可以可加"-"前缀，表示减去对应的时间
func CustomStringToSecond(ctx tyto.Context, baseTime time.Time, str string) (int64, bool) {
	var (
		month int64
		week  int64
		day   int64
		hour  int64
		min   int64
		sec   int64
	)

	l := strings.Split(str, " ")
	for _, v := range l {
		if len(v) < 2 {
			ctx.Logger().Error("invalid string:", v)
			return 0, false
		}

		num, err := strconv.ParseInt(v[:len(v)-1], 10, 64)
		if err != nil {
			ctx.Logger().Error("parse int failed:", err)
			return 0, false
		}

		switch v[len(v)-1] {
		case 'M':
			month = num
		case 'w':
			week = num
		case 'd':
			day = num
		case 'h':
			hour = num
		case 'm':
			min = num
		case 's':
			sec = num
		default:
			ctx.Logger().Error("invalid string:", v)
			return 0, false
		}

	}

	total := int64(0)
	if month > 0 {
		begin := BeginningOfMonth(baseTime)
		end := begin.AddDate(0, int(month), 0)
		total += end.Unix() - begin.Unix()
	}

	total += week * TOTAL_SECONDS_PER_WEEK
	total += day * TOTAL_SECONDS_PER_DAY
	total += hour * TOTAL_SECONDS_PER_HOUR
	total += min * TOTAL_SECONDS_PER_MINUTE
	total += sec

	return total, true
}
