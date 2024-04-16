package internal

import "os"

func NewWriter(file *os.File, bufferSize int32) Writer {
	if bufferSize > 0 {
		return NewBufferedWriter(file, bufferSize)
	}

	return NewUnbufferedWriter(file)
}
