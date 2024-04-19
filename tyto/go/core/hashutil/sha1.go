package hashutil

import (
	"crypto/sha1"
	"tyto/core/tyto"
)

// 计算字符串的sha1值
func Sha1Sum(ctx tyto.Context, format string, a ...interface{}) string {
	h := sha1.New()
	return hashSum(ctx, h, format, a...)
}

func Sha1SumToBase64(ctx tyto.Context, format string, a ...interface{}) string {
	h := sha1.New()
	return hashSumToBase64(ctx, h, format, a...)
}
