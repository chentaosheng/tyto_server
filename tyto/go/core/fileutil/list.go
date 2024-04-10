package fileutil

import (
	"path"
	"path/filepath"
)

// 列出指定目录下的文件
func ListFile(dir string, globPattern string) ([]string, error) {
	p := path.Join(dir, globPattern)
	return filepath.Glob(p)
}
