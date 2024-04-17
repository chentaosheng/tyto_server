package logs

import "bytes"

// 单个日志项格式化器
type Formatter interface {
	// 格式化日志
	Format(buff *bytes.Buffer, record Record) error
}
