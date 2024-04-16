package internal

import (
	"github.com/itchyny/timefmt-go"
	"regexp"
	"time"
)

func GetBaseTime(t time.Time, interval time.Duration) time.Time {
	if t.Location() != time.UTC {
		// 转换为UTC时间
		year, month, day := t.Date()
		hour, min, sec := t.Clock()
		nsec := t.Nanosecond()
		baseTime := time.Date(year, month, day, hour, min, sec, nsec, time.UTC)

		// 非utc时间进行Truncate操作会带上时区信息，所以需要手动去掉
		baseTime = baseTime.Truncate(interval)

		// 还原为原始时区
		year, month, day = baseTime.Date()
		hour, min, sec = baseTime.Clock()
		nsec = baseTime.Nanosecond()

		return time.Date(year, month, day, hour, min, sec, nsec, t.Location())
	}

	return t.Truncate(interval)
}

func GenerateFileName(format string, t time.Time) string {
	return timefmt.Format(t, format)
}

func ToGlobPattern(pattern string) string {
	// 转换表
	var conversionList = []*regexp.Regexp{
		// 替换%*为*
		regexp.MustCompile(`%[%+A-Za-z]`),
		// 替换**为*
		regexp.MustCompile(`\*+`),
	}

	globPattern := pattern
	for _, re := range conversionList {
		globPattern = re.ReplaceAllString(globPattern, "*")
	}

	return globPattern
}
