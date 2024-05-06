package logs

// 日志级别
type Level int32

const (
	LEVEL_TRACE Level = 0 // 跟踪
	LEVEL_DEBUG Level = 1 // 调试
	LEVEL_INFO  Level = 2 // 信息
	LEVEL_WARN  Level = 3 // 警告
	LEVEL_ERROR Level = 4 // 错误
	LEVEL_FATAL Level = 5 // 致命
)

const (
	LEVEL_MIN     Level = LEVEL_TRACE   // 最小级别
	LEVEL_MAX     Level = LEVEL_FATAL   // 最大级别
	LEVEL_INVALID Level = LEVEL_MAX + 1 // 无效级别
)

// 用于将级别转换为名称
var sLevelNameList = [LEVEL_MAX + 2][]byte{
	LEVEL_TRACE:   []byte("TRACE"),
	LEVEL_DEBUG:   []byte("DEBUG"),
	LEVEL_INFO:    []byte("INFO"),
	LEVEL_WARN:    []byte("WARN"),
	LEVEL_ERROR:   []byte("ERROR"),
	LEVEL_FATAL:   []byte("FATAL"),
	LEVEL_INVALID: []byte("UNKNOWN"),
}

// 将级别转换为名称
func (level Level) Marshal() []byte {
	if level < LEVEL_MIN || level > LEVEL_MAX {
		return sLevelNameList[LEVEL_INVALID]
	}

	return sLevelNameList[level]
}
