package rolling

import (
	"time"
	"tyto/core/logging"
	"tyto/core/rolling/internal"
)

// 单个选项
type Option func(*internal.Options)

// 文件输出目录
func WithOutDir(outDir string) Option {
	return func(o *internal.Options) {
		o.OutDir = outDir
	}
}

// namePattern参数，用于文件名生成，"app.txt.%F" 会生成 "app.txt.2006-01-02"，
// 规则:
//
//	%F 年-月-日，"2006-01-02"
//	%Y 年，"2006"
//	%m 月，"01"
//	%d 日，"02"
//	%H 小时，"15"
//	%M 分钟，"04"
func WithNamePattern(namePattern string) Option {
	return func(o *internal.Options) {
		o.NamePattern = namePattern
	}
}

// linkName参数只在linux下有效，输入空字符串表示不创建符号链接
func WithLinkName(linkName string) Option {
	return func(o *internal.Options) {
		o.LinkName = linkName
	}
}

// 文件保留时间，<=0表示不清理
func WithMaxAge(maxAge time.Duration) Option {
	return func(o *internal.Options) {
		o.MaxAge = maxAge
	}
}

// 清理过期文件的时间间隔
func WithCleanupInterval(cleanupInterval time.Duration) Option {
	return func(o *internal.Options) {
		o.CleanupInterval = cleanupInterval
	}
}

// 文件轮转间隔
func WithRotationInterval(rotationInterval time.Duration) Option {
	return func(o *internal.Options) {
		o.RotationInterval = rotationInterval
	}
}

// 写入队列大小，会自动扩容
func WithWriterQueueSize(writerQueueSize int32) Option {
	return func(o *internal.Options) {
		o.WriterQueueSize = writerQueueSize
	}
}

// 写入缓冲区大小
func WithBufferSize(bufferSize int32) Option {
	return func(o *internal.Options) {
		o.BufferSize = bufferSize
	}
}

// 缓冲区刷到文件的时间间隔
func WithFlushInterval(flushInterval time.Duration) Option {
	return func(o *internal.Options) {
		o.FlushInterval = flushInterval
	}
}

// 日志器
func WithLogger(logger logging.Logger) Option {
	return func(o *internal.Options) {
		o.Logger = logger
	}
}
