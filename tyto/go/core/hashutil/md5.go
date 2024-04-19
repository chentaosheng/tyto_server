package hashutil

import (
	"crypto/md5"
	"tyto/core/tyto"
)

// 计算字符串的md5值
func Md5Sum(ctx tyto.Context, format string, a ...interface{}) string {
	h := md5.New()
	return hashSum(ctx, h, format, a...)
}

func Md5SumToBase64(ctx tyto.Context, format string, a ...interface{}) string {
	h := md5.New()
	return hashSumToBase64(ctx, h, format, a...)
}
