package logs

import (
	"bytes"
	"io"
	"sync"
	"tyto/core/logs/mini"
)

// 用于性能测试的日志接收器
type DiscardSink struct {
	logger    *mini.Logger
	formatter Formatter
	pool      sync.Pool
}

func NewDiscardSink(logger *mini.Logger) Sink {
	return &DiscardSink{
		formatter: NewTextFormatter(DEFAULT_SKIP_CALLER_COUNT, DEFAULT_MAX_CALLER_COUNT),
		pool: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 0, 128))
			},
		},
	}
}

func (sink *DiscardSink) Handle(record Record) {
	buff := sink.pool.Get().(*bytes.Buffer)
	defer sink.pool.Put(buff)

	buff.Reset()

	if err := sink.formatter.Format(buff, record); err != nil {
		sink.logger.Error("failed to format log record, err:", err.Error())
		return
	}

	data := buff.Bytes()
	_, err := io.Discard.Write(data)
	if err != nil {
		sink.logger.Error("failed to write log record, err:", err.Error())
		return
	}
}

func (sink *DiscardSink) Close() {
}
