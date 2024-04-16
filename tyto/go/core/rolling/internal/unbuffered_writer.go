package internal

import "os"

// 无缓冲
type UnbufferedWriter struct {
	file *os.File
}

func NewUnbufferedWriter(file *os.File) Writer {
	return &UnbufferedWriter{
		file: file,
	}
}

func (writer *UnbufferedWriter) Write(p []byte) (n int, err error) {
	return writer.file.Write(p)
}

func (writer *UnbufferedWriter) WriteString(s string) (n int, err error) {
	return writer.file.WriteString(s)
}

func (writer *UnbufferedWriter) IsBuffered() bool {
	return false
}

func (writer *UnbufferedWriter) Sync() error {
	return writer.file.Sync()
}

func (writer *UnbufferedWriter) Close() error {
	writer.file.Sync()
	return writer.file.Close()
}

func (writer *UnbufferedWriter) Reset(file *os.File) {
	writer.file.Sync()
	writer.file.Close()
	writer.file = file
}

// file对象指向的文件是否已经被删除
func (writer *UnbufferedWriter) IsValid() bool {
	_, err := writer.file.Stat()
	return err == nil
}
