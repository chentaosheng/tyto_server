package logs

import (
	"bytes"
	"sync"
	"sync/atomic"
	"time"
	"tyto/core/logs/mini"
	"tyto/core/memutil"
	"tyto/core/rolling"
)

// 文件日志接收器
type FileSink struct {
	logger       *mini.Logger
	formatter    Formatter
	pool         sync.Pool
	normalWriter *rolling.RotateWriter
	errorWriter  *rolling.RotateWriter
	closed       atomic.Bool
}

func NewDefaultFileSink(logger *mini.Logger, outDir string, logFileName string) Sink {
	formatter := NewTextFormatter(DEFAULT_SKIP_CALLER_COUNT, DEFAULT_MAX_CALLER_COUNT)
	return NewFileSink(logger, outDir, logFileName, formatter)
}

func NewFileSink(logger *mini.Logger, outDir string, logFileName string, formatter Formatter) Sink {
	// 普通日志
	normalWriter, err := newFileSinkWriter(logger, outDir, logFileName, 2048, true)
	if err != nil {
		logger.Error("failed to create normal file writer, err:", err.Error())
		return nil
	}

	// 错误日志
	errorWriter, err := newFileSinkWriter(logger, outDir, "err_"+logFileName, 256, true)
	if err != nil {
		logger.Error("failed to create error file writer, err:", err.Error())
		return nil
	}

	// create
	sink := &FileSink{
		logger:       logger,
		formatter:    formatter,
		pool:         sync.Pool{New: nil},
		normalWriter: normalWriter,
		errorWriter:  errorWriter,
		closed:       atomic.Bool{},
	}

	sink.pool.New = func() interface{} {
		buff := bytes.Buffer{}
		buff.Grow(128)
		return memutil.NewRefObject(buff, sink.destroyBuffer)
	}

	return sink
}

// 创建
func newFileSinkWriter(logger *mini.Logger, outDir string, logFileName string, queueSize int32, buffered bool) (*rolling.RotateWriter, error) {
	var (
		bufferSize    int32
		flushInterval time.Duration
	)

	if buffered {
		// 有缓冲模式
		bufferSize = 128 * 1024
		flushInterval = 5 * time.Second

	} else {
		// 无缓冲模式
		bufferSize = 0
		flushInterval = 0
	}

	return rolling.NewRotateWriter(
		rolling.WithOutDir(outDir),
		rolling.WithNamePattern(logFileName+".%F"),
		rolling.WithLinkName(logFileName),
		rolling.WithMaxAge(14*24*time.Hour),
		rolling.WithCleanupInterval(24*time.Hour),
		rolling.WithRotationInterval(24*time.Hour),
		rolling.WithWriterQueueSize(queueSize),
		rolling.WithBufferSize(bufferSize),
		rolling.WithFlushInterval(flushInterval),
		rolling.WithLogger(logger),
	)
}

func (sink *FileSink) destroyBuffer(buff *rolling.BufferType) {
	sink.pool.Put(buff)
}

func (sink *FileSink) Handle(record Record) {
	if sink.closed.Load() {
		return
	}

	buff := sink.pool.Get().(*rolling.BufferType)
	buff.IncRef()
	buff.Object().Reset()
	defer buff.DecRef()

	if err := sink.formatter.Format(buff.Object(), record); err != nil {
		sink.logger.Error("failed to format log record, err:", err.Error())
		return
	}

	// 普通日志
	if _, err := sink.normalWriter.WriteBuffer(buff); err != nil {
		sink.logger.Error("failed to write normal log record, err:", err.Error())
	}

	// 错入日志
	if record.GetLevel() >= LEVEL_WARN {
		if _, err := sink.errorWriter.WriteBuffer(buff); err != nil {
			sink.logger.Error("failed to write error log record, err:", err.Error())
		}
	}
}

func (sink *FileSink) Close() {
	if sink.closed.Load() {
		return
	}
	if !sink.closed.CompareAndSwap(false, true) {
		return
	}

	if err := sink.normalWriter.Close(); err != nil {
		sink.logger.Error("failed to close normal file writer, err:", err.Error())
	}

	if err := sink.errorWriter.Close(); err != nil {
		sink.logger.Error("failed to close error file writer, err:", err.Error())
	}
}
