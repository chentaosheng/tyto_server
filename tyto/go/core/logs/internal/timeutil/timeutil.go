package timeutil

import (
	"bytes"
	"time"
)

func appendInt(buff *bytes.Buffer, x int, width int) {
	u := uint(x)
	if x < 0 {
		buff.WriteByte('-')
		u = uint(-x)
	}

	var buf [20]byte
	i := len(buf)
	for u >= 10 {
		i--
		q := u / 10
		buf[i] = byte('0' + u - q*10)
		u = q
	}
	i--
	buf[i] = byte('0' + u)

	// 补0
	for w := len(buf) - i; w < width; w++ {
		buff.WriteByte('0')
	}

	buff.Write(buf[i:])
}

// 将时间格式化为19-01-01 00:00:00.000
func FormatShort(buff *bytes.Buffer, t time.Time) {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	ms := t.Nanosecond() / int(time.Millisecond)

	// year
	appendInt(buff, year%100, 2)

	// month
	buff.WriteByte('-')
	appendInt(buff, int(month), 2)

	// day
	buff.WriteByte('-')
	appendInt(buff, day, 2)

	buff.WriteByte(' ')

	// hour
	appendInt(buff, hour, 2)

	// min
	buff.WriteByte(':')
	appendInt(buff, min, 2)

	// sec
	buff.WriteByte(':')
	appendInt(buff, sec, 2)

	// ms
	buff.WriteByte('.')
	appendInt(buff, ms, 3)
}

// 将t格式化为2019-01-01 00:00:00.000
func Format(buff *bytes.Buffer, t time.Time) {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	ms := t.Nanosecond() / int(time.Millisecond)

	// year
	appendInt(buff, year, 4)

	// month
	buff.WriteByte('-')
	appendInt(buff, int(month), 2)

	// day
	buff.WriteByte('-')
	appendInt(buff, day, 2)

	buff.WriteByte(' ')

	// hour
	appendInt(buff, hour, 2)

	// min
	buff.WriteByte(':')
	appendInt(buff, min, 2)

	// sec
	buff.WriteByte(':')
	appendInt(buff, sec, 2)

	// ms
	buff.WriteByte('.')
	appendInt(buff, ms, 3)
}
