package logs

import (
	"time"
)

// 文本型日志记录
type TextRecord struct {
	RecordType       RecordType
	Level            Level
	ReportCallerType ReportCallerType
	Time             time.Time
	Args             []interface{}
}

func NewTextRecord() *TextRecord {
	return &TextRecord{
		RecordType: RECORD_TYPE_TEXT,
	}
}

func (r *TextRecord) Reset(level Level, reportCallerType ReportCallerType, args []interface{}) {
	r.RecordType = RECORD_TYPE_TEXT
	r.Level = level
	r.ReportCallerType = reportCallerType
	r.Time = time.Now().Local()
	r.Args = args
}

func (r *TextRecord) GetRecordType() RecordType {
	return r.RecordType
}

func (r *TextRecord) GetLevel() Level {
	return r.Level
}

func (r *TextRecord) GetReportCallerType() ReportCallerType {
	return r.ReportCallerType
}
