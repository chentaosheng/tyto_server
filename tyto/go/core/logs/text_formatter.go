package logs

import (
	"bytes"
	"fmt"
	"tyto/core/logs/internal/stackutil"
	"tyto/core/logs/internal/timeutil"
)

const (
	DEFAULT_SKIP_CALLER_COUNT int32 = 7 // 默认调用栈跳过的层数
	DEFAULT_MAX_CALLER_COUNT  int32 = 8 // 默认打印调用栈的最大层数
)

type TextFormatter struct {
	skipCallerCount int32
	maxCallerCount  int32
}

func NewTextFormatter(skip, max int32) *TextFormatter {
	return &TextFormatter{
		skipCallerCount: skip,
		maxCallerCount:  max,
	}
}

func (f *TextFormatter) Format(buff *bytes.Buffer, record Record) error {
	switch record.GetRecordType() {
	case RECORD_TYPE_TEXT:
		return f.formatTextRecord(buff, record)
	default:
		return fmt.Errorf("record type not supported, type: %d", record.GetRecordType())
	}
}

func (f *TextFormatter) formatTextRecord(buff *bytes.Buffer, record Record) error {
	r := record.(*TextRecord)

	// 时间
	timeutil.Format(buff, r.Time)

	// 日志级别
	buff.WriteByte(' ')
	buff.WriteByte('[')
	buff.Write(r.Level.Marshal())
	buff.WriteByte(']')

	// 日志内容
	buff.WriteByte(' ')
	fmt.Fprintln(buff, r.Args...)

	// 调用栈
	return f.reportCallers(buff, record)
}

func (f *TextFormatter) reportCallers(buff *bytes.Buffer, record Record) error {
	switch record.GetReportCallerType() {
	case REPORT_CALLER_TYPE_NONE:
		return nil

	case REPORT_CALLER_TYPE_ERROR:
		if record.GetLevel() >= LEVEL_ERROR {
			stackutil.Print(buff, f.skipCallerCount, f.maxCallerCount)
		}

	case REPORT_CALLER_TYPE_ALWAYS:
		stackutil.Print(buff, f.skipCallerCount, f.maxCallerCount)

	default:
		return fmt.Errorf("report caller type not supported, type: %d", record.GetReportCallerType())
	}

	return nil
}
