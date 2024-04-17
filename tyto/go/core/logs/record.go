package logs

// 日志记录类型
type RecordType int32

const (
	RECORD_TYPE_TEXT RecordType = 0 // 文本
)

// 打印调用栈的类型
type ReportCallerType int32

const (
	REPORT_CALLER_TYPE_NONE   ReportCallerType = 0 // 任何时候都不打印
	REPORT_CALLER_TYPE_ERROR  ReportCallerType = 1 // 错误级别以上打印
	REPORT_CALLER_TYPE_ALWAYS ReportCallerType = 2 // 所有级别都打印
)

// 日志记录对象
type Record interface {
	// 获取日志记录类型
	GetRecordType() RecordType
	// 获取日志级别
	GetLevel() Level
	// 堆栈跟踪类型
	GetReportCallerType() ReportCallerType
}
