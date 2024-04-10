package fileutil

import (
	"strings"
)

// 转换为Windows风格的路径分隔符
func ToWindowsSeparator(path string) string {
	return strings.ReplaceAll(path, "/", "\\")
}

// 转换为Linux风格的路径分隔符
func ToLinuxSeparator(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}
