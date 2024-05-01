package timeutil

// 时间常量
const (
	TOTAL_SECONDS_PER_MINUTE      int64 = 60                      // 1分钟总秒数
	TOTAL_SECONDS_PER_HOUR        int64 = 60 * 60                 // 1小时总秒数
	TOTAL_SECONDS_PER_DAY         int64 = 24 * 60 * 60            // 1天总秒数
	TOTAL_SECONDS_PER_WEEK        int64 = 7 * 24 * 60 * 60        // 1周总秒数
	TOTAL_MILLISECONDS_PER_SECOND int64 = 1000                    // 1秒总毫秒数
	TOTAL_MILLISECONDS_PER_MINUTE int64 = 60 * 1000               // 1分钟总毫秒数
	TOTAL_MILLISECONDS_PER_HOUR   int64 = 60 * 60 * 1000          // 1小时总毫秒数
	TOTAL_MILLISECONDS_PER_DAY    int64 = 24 * 60 * 60 * 1000     // 1天总毫秒数
	TOTAL_MILLISECONDS_PER_WEEK   int64 = 7 * 24 * 60 * 60 * 1000 // 1周总毫秒数
)
