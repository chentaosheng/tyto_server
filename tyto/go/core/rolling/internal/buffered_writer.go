package internal

import (
	"bufio"
	"os"
)

// 有缓冲
type BufferedWriter struct {
	file *os.File      // 当前的输出文件
	buff *bufio.Writer // buffer io时，使用的writer
}

func NewBufferedWriter(file *os.File, bufferSize int32) Writer {
	return &BufferedWriter{
		file: file,
		buff: bufio.NewWriterSize(file, int(bufferSize)),
	}
}

func (writer *BufferedWriter) Write(p []byte) (n int, err error) {
	return writer.buff.Write(p)
}

func (writer *BufferedWriter) WriteString(s string) (n int, err error) {
	return writer.buff.WriteString(s)
}

func (writer *BufferedWriter) IsBuffered() bool {
	return true
}

func (writer *BufferedWriter) Sync() error {
	writer.buff.Flush()
	return writer.file.Sync()
}

func (writer *BufferedWriter) Close() error {
	writer.Sync()
	return writer.file.Close()
}

func (writer *BufferedWriter) Reset(file *os.File) {
	writer.buff.Flush()
	writer.buff.Reset(file)

	writer.file.Sync()
	writer.file.Close()
	writer.file = file
}

// file对象指向的文件是否已经被删除
func (writer *BufferedWriter) IsValid() bool {
	_, err := writer.file.Stat()
	return err == nil
}
