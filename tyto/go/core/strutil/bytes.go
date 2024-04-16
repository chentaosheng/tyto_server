package strutil

import "unsafe"

// bytes转换为string，不会发生内存拷贝
func BytesToStr(bs []byte) string {
	if len(bs) == 0 {
		return ""
	}

	return unsafe.String(&bs[0], len(bs))
}

// string转换为bytes，不会发生内存拷贝
func StrToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
