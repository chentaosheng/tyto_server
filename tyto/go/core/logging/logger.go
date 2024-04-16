package logging

// logger接口
type Logger interface {
	// 跟踪
	Trace(v ...interface{})
	// 调试
	Debug(v ...interface{})
	// 信息
	Info(v ...interface{})
	// 警告
	Warn(v ...interface{})
	// 错误
	Error(v ...interface{})
	// 致命
	Fatal(v ...interface{})

	// 关闭日志系统
	Close()
}
