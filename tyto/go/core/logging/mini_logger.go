package logging

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// 一个简单的日志器，只能输出到stderr
// 通常用于在日志系统初始化之前使用
type MiniLogger struct {
	mutex   sync.Mutex
	builder strings.Builder
}

func NewMiniLogger() *MiniLogger {
	logger := &MiniLogger{
		mutex:   sync.Mutex{},
		builder: strings.Builder{},
	}

	logger.builder.Grow(256)

	return logger
}

func (logger *MiniLogger) Trace(v ...interface{}) {
	logger.log("TRACE", v...)
}

func (logger *MiniLogger) Debug(v ...interface{}) {
	logger.log("DEBUG", v...)
}

func (logger *MiniLogger) Info(v ...interface{}) {
	logger.log("INFO", v...)
}

func (logger *MiniLogger) Warn(v ...interface{}) {
	logger.log("WARN", v...)
}

func (logger *MiniLogger) Error(v ...interface{}) {
	logger.log("ERROR", v...)
}

func (logger *MiniLogger) Fatal(v ...interface{}) {
	logger.log("FATAL", v...)
}

func (logger *MiniLogger) Close() {
}

func (logger *MiniLogger) log(level string, v ...interface{}) {
	logger.mutex.Lock()
	defer logger.mutex.Unlock()

	logger.builder.Reset()

	// 日期
	t := time.Now().Local().Format("2006-01-02 15:04:05.000")
	logger.builder.WriteString(t)

	// 日志级别
	logger.builder.WriteString(" [")
	logger.builder.WriteString(level)
	logger.builder.WriteString("] ")

	// 日志内容
	s := fmt.Sprintln(v...)
	logger.builder.WriteString(s)

	// 打印
	fmt.Fprintf(os.Stderr, logger.builder.String())
}
