package osutil

import (
	"runtime"
	"unsafe"
)

// 获取软件位数
func GetBit() int32 {
	return int32(unsafe.Sizeof(uintptr(0)) * 8)
}

// 是否64位软件
func Is64Bit() bool {
	return GetBit() == 64
}

// 是否32为软件
func Is32Bit() bool {
	return GetBit() == 32
}

// 是否windows系统
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// 是否linux系统
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// 是否mac osx系统
func IsMacOsx() bool {
	return runtime.GOOS == "darwin"
}
