package fileutil

import (
	"path/filepath"
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

// 移除路径末尾的分隔符
func RemoveLastSeparator(path string) string {
	if strings.HasSuffix(path, "/") || strings.HasSuffix(path, "\\") {
		return path[:len(path)-1]
	}

	return path
}

// 添加路径末尾的分隔符
func AppendLastSeparator(path string) string {
	if strings.HasSuffix(path, "/") || strings.HasSuffix(path, "\\") {
		return path
	}

	b := strings.Builder{}
	b.Grow(len(path) + 1)

	b.WriteString(path)
	b.WriteByte(filepath.Separator)

	return b.String()
}
