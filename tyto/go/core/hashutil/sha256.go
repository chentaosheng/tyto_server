package hashutil

import (
	"crypto/sha256"
	"tyto/core/tyto"
)

// 计算字符串的sha256值
func Sha256Sum(ctx tyto.Context, format string, a ...interface{}) string {
	h := sha256.New()
	return hashSum(ctx, h, format, a...)
}

func Sha256SumToBase64(ctx tyto.Context, format string, a ...interface{}) string {
	h := sha256.New()
	return hashSumToBase64(ctx, h, format, a...)
}
