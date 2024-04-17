package logs

// 日志级别
type Level int32

const (
	LEVEL_TRACE Level = 1 // 跟踪
	LEVEL_DEBUG Level = 2 // 调试
	LEVEL_INFO  Level = 3 // 信息
	LEVEL_WARN  Level = 4 // 警告
	LEVEL_ERROR Level = 5 // 错误
	LEVEL_FATAL Level = 6 // 致命
)

const (
	MIN_LEVEL     Level = LEVEL_TRACE   // 最小级别
	MAX_LEVEL     Level = LEVEL_FATAL   // 最大级别
	INVALID_LEVEL Level = MIN_LEVEL - 1 // 无效级别
)

var (
	// 未知级别名称
	sUnknownLevelName = []byte("UNKNOWN")

	// 用于将级别转换为名称
	// 第0个元素无效，使用level作为下标
	sLevelNameList = [MAX_LEVEL + 1][]byte{
		INVALID_LEVEL: sUnknownLevelName,
		LEVEL_TRACE:   []byte("TRACE"),
		LEVEL_DEBUG:   []byte("DEBUG"),
		LEVEL_INFO:    []byte("INFO"),
		LEVEL_WARN:    []byte("WARN"),
		LEVEL_ERROR:   []byte("ERROR"),
		LEVEL_FATAL:   []byte("FATAL"),
	}
)

// 将级别转换为名称
func (level Level) Marshal() []byte {
	if level < MIN_LEVEL || level > MAX_LEVEL {
		return sUnknownLevelName
	}

	return sLevelNameList[level]
}
