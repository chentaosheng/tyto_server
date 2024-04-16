package internal

import (
	"io"
	"os"
)

// 文件写入接口
type Writer interface {
	io.WriteCloser
	io.StringWriter

	// 是否有缓冲
	IsBuffered() bool
	// 将缓冲区的数据写入文件
	Sync() error
	// 重置文件
	Reset(file *os.File)
	// 文件是否有效
	IsValid() bool
}
