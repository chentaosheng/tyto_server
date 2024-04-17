package stackutil

import (
	"bytes"
	"runtime"
	"strconv"
)

// 打印函数调用栈信息到buff中
// skip: 跳过栈信息的层级，从开始部分开始计算
// max: 最大打印多个函数信息
func Print(buff *bytes.Buffer, skip int32, max int32) {
	pcs := make([]uintptr, max)
	count := runtime.Callers(int(skip+1), pcs)

	buff.WriteString("stack:\n")

	for i := 0; i < count; i++ {
		pc := pcs[i] - 1
		formatFrame(buff, pc)
	}
}

func formatFrame(buff *bytes.Buffer, pc uintptr) {
	f := runtime.FuncForPC(pc)
	if f == nil {
		buff.WriteString("unknown")
		return
	}

	buff.WriteByte('\t')
	buff.WriteString(f.Name())
	buff.WriteString("()\n\t\t")

	file, line := f.FileLine(pc)
	buff.WriteString(file)
	buff.WriteByte(':')
	buff.WriteString(strconv.Itoa(line))
	buff.WriteByte('\n')
}
