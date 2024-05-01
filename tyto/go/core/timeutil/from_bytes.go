package timeutil

import (
	"time"
	"tyto/core/byteutil"
	"tyto/core/tyto"
)

// 将yyyy-MM-dd HH:mm:ss格式的字节数组转换为time.Time
func BytesToTime(ctx tyto.Context, bs []byte) (time.Time, bool) {
	str := byteutil.BytesToStr(bs)
	t, err := time.ParseInLocation(time.DateTime, str, time.Local)
	if err != nil {
		ctx.Logger().Error("convert bytes to time failed:", err)
		return t, false
	}

	return t, true
}

// 将yyyy-MM-dd HH:mm:ss格式的字节数组转换为1970开始的毫秒数
func BytesToMilli(ctx tyto.Context, bs []byte) (int64, bool) {
	t, ok := BytesToTime(ctx, bs)
	if !ok {
		return 0, false
	}
	return GetMilli(t), true
}

// 将yyyy-MM-dd HH:mm:ss格式的字符串转换为1970开始的秒数
func BytesToSecond(ctx tyto.Context, bs []byte) (int64, bool) {
	t, ok := BytesToTime(ctx, bs)
	if !ok {
		return 0, false
	}
	return GetSecond(t), true
}
