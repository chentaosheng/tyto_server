package logs

import (
	"bytes"
	"os"
	"sync"
	"tyto/core/logs/mini"

	colortext "github.com/daviddengcn/go-colortext"
)

// 输出到控制台的日志接收器
// 不要用于生产环境，性能较差
type ConsoleSink struct {
	logger    *mini.Logger
	formatter Formatter
	mutex     sync.Mutex
	pool      sync.Pool
}

func NewDefaultConsoleSink(logger *mini.Logger) Sink {
	formatter := NewTextFormatter(DEFAULT_SKIP_CALLER_COUNT, DEFAULT_MAX_CALLER_COUNT)
	return NewConsoleSink(logger, formatter)
}

func NewConsoleSink(logger *mini.Logger, formatter Formatter) Sink {
	return &ConsoleSink{
		logger:    logger,
		formatter: formatter,
		mutex:     sync.Mutex{},
		pool: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, 128))
			},
		},
	}
}

func (sink *ConsoleSink) Handle(record Record) {
	buff := sink.pool.Get().(*bytes.Buffer)
	defer sink.pool.Put(buff)

	buff.Reset()

	if err := sink.formatter.Format(buff, record); err != nil {
		sink.logger.Error("failed to format log record, err:", err.Error())
		return
	}

	data := buff.Bytes()

	sink.mutex.Lock()
	defer sink.mutex.Unlock()

	sink.BeginPaint(record)
	defer sink.EndPaint(record)

	_, err := os.Stdout.Write(data)
	if err != nil {
		sink.logger.Error("failed to write log record, err:", err.Error())
		return
	}
}

func (sink *ConsoleSink) Close() {
	sink.mutex.Lock()
	defer sink.mutex.Unlock()

	os.Stdout.Sync()
}

func (sink *ConsoleSink) BeginPaint(record Record) {
	switch record.GetLevel() {
	case LEVEL_WARN:
		colortext.Foreground(colortext.Yellow, true)
	case LEVEL_ERROR:
		colortext.Foreground(colortext.Red, true)
	case LEVEL_FATAL:
		colortext.Foreground(colortext.Magenta, true)
	}
}

func (sink *ConsoleSink) EndPaint(record Record) {
	colortext.ResetColor()
}
