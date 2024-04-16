package internal

import (
	"errors"
	"path/filepath"
	"strings"
	"time"
	"tyto/core/logging"
)

const (
	DEFAULT_MAX_AGE           = 30 * 24 * time.Hour // 默认文件保留时间
	DEFAULT_CLEANUP_INTERVAL  = 24 * time.Hour      // 默认文件清理间隔
	DEFAULT_ROTATION_INTERVAL = 24 * time.Hour      // 默认文件轮转间隔
	DEFAULT_FLUSH_INTERVAL    = 5 * time.Second     // 默认缓冲区刷到文件的时间间隔
	DEFAULT_WRITER_QUEUE_SIZE = 1024                // 默认写入队列大小
	DEFAULT_BUFFER_SIZE       = 64 * 1024           // 默认缓冲区大小
)

// 旋转文件选项
type Options struct {
	OutDir           string         // 文件输出目录
	NamePattern      string         // 文件名生成模式
	LinkName         string         // 文件符号链接名，用于给当前正在写入的文件创建一个符号链接
	MaxAge           time.Duration  // 文件保留时间，<=0表示不清理
	CleanupInterval  time.Duration  // 清理过期文件的时间间隔
	RotationInterval time.Duration  // 文件轮转间隔
	WriterQueueSize  int32          // 写入队列大小，会自动扩容
	BufferSize       int32          // 写入缓冲区大小
	FlushInterval    time.Duration  // 缓冲区刷到文件的时间间隔
	Logger           logging.Logger // 日志器
}

// 新建旋转文件选项
func NewDefaultOptions() *Options {
	return &Options{
		OutDir:           "",
		NamePattern:      "",
		LinkName:         "",
		MaxAge:           DEFAULT_MAX_AGE,
		CleanupInterval:  DEFAULT_CLEANUP_INTERVAL,
		RotationInterval: DEFAULT_ROTATION_INTERVAL,
		WriterQueueSize:  DEFAULT_WRITER_QUEUE_SIZE,
		BufferSize:       DEFAULT_BUFFER_SIZE,
		FlushInterval:    DEFAULT_FLUSH_INTERVAL,
		Logger:           nil,
	}
}

// 检查选项合法性
func (o *Options) Validate() error {
	if o.OutDir == "" {
		return errors.New("OutDir is required")
	}

	_, err := filepath.Abs(o.OutDir)
	if err != nil {
		return err
	}

	if o.NamePattern == "" {
		return errors.New("NamePattern is required")
	}

	if strings.ContainsAny(o.NamePattern, "*$/\\") {
		return errors.New("NamePattern format error")
	}

	if o.CleanupInterval <= 0 {
		return errors.New("CleanupInterval must be greater than 0")
	}

	if o.RotationInterval <= 0 {
		return errors.New("RotationInterval must be greater than 0")
	}

	if o.BufferSize < 0 {
		return errors.New("BufferSize must be greater than or equal to 0")
	}

	if o.BufferSize > 0 && o.FlushInterval <= 0 {
		return errors.New("FlushInterval must be greater than 0 when BufferSize is greater than 0")
	}

	if o.Logger == nil {
		return errors.New("Logger is required")
	}

	return nil
}
