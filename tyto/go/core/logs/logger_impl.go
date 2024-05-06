package logs

import (
	"sync"
	"sync/atomic"
)

type LoggerImpl struct {
	pool  sync.Pool
	level atomic.Int32
	sinks []Sink
}

func NewLoggerImpl(sinks ...Sink) *LoggerImpl {
	logger := &LoggerImpl{
		pool: sync.Pool{
			New: func() interface{} {
				return NewTextRecord()
			},
		},
		level: atomic.Int32{},
		sinks: sinks,
	}

	logger.level.Store(int32(LEVEL_DEBUG))
	return logger
}

func (logger *LoggerImpl) IsEnabled(level Level) bool {
	return level >= Level(logger.GetLevel())
}

func (logger *LoggerImpl) Trace(v ...interface{}) {
	if !logger.IsEnabled(LEVEL_TRACE) {
		return
	}
	logger.log(LEVEL_TRACE, REPORT_CALLER_TYPE_ERROR, v)
}

func (logger *LoggerImpl) Debug(v ...interface{}) {
	if !logger.IsEnabled(LEVEL_DEBUG) {
		return
	}
	logger.log(LEVEL_DEBUG, REPORT_CALLER_TYPE_ERROR, v)
}

func (logger *LoggerImpl) Info(v ...interface{}) {
	if !logger.IsEnabled(LEVEL_INFO) {
		return
	}
	logger.log(LEVEL_INFO, REPORT_CALLER_TYPE_ERROR, v)
}

func (logger *LoggerImpl) Warn(v ...interface{}) {
	if !logger.IsEnabled(LEVEL_WARN) {
		return
	}
	logger.log(LEVEL_WARN, REPORT_CALLER_TYPE_ERROR, v)
}

func (logger *LoggerImpl) Error(v ...interface{}) {
	if !logger.IsEnabled(LEVEL_ERROR) {
		return
	}
	logger.log(LEVEL_ERROR, REPORT_CALLER_TYPE_ERROR, v)
}

func (logger *LoggerImpl) Fatal(v ...interface{}) {
	if !logger.IsEnabled(LEVEL_FATAL) {
		return
	}
	logger.log(LEVEL_FATAL, REPORT_CALLER_TYPE_ERROR, v)
}

func (logger *LoggerImpl) NoCallerError(v ...interface{}) {
	if !logger.IsEnabled(LEVEL_ERROR) {
		return
	}
	logger.log(LEVEL_ERROR, REPORT_CALLER_TYPE_NONE, v)
}

func (logger *LoggerImpl) NoCallerFatal(v ...interface{}) {
	if !logger.IsEnabled(LEVEL_FATAL) {
		return
	}
	logger.log(LEVEL_FATAL, REPORT_CALLER_TYPE_NONE, v)
}

func (logger *LoggerImpl) GetLevel() int32 {
	return logger.level.Load()
}

// 设置日志过滤级别，只有大于等于该级别的日志才会被输出
// 级别错误会导致panic
func (logger *LoggerImpl) SetLevel(level int32) {
	if level < int32(LEVEL_MIN) || level > int32(LEVEL_MAX) {
		panic("invalid log level")
	}

	logger.level.Store(level)
}

func (logger *LoggerImpl) Close() {
	for _, sink := range logger.sinks {
		if sink == nil {
			continue
		}

		sink.Close()
	}
}

func (logger *LoggerImpl) log(level Level, reportCallerType ReportCallerType, v []interface{}) {
	record := logger.pool.Get().(*TextRecord)
	defer logger.pool.Put(record)

	record.Reset(level, reportCallerType, v)

	for _, sink := range logger.sinks {
		if sink == nil {
			continue
		}

		sink.Handle(record)
	}
}
