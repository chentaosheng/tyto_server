package hashutil

import (
	"crypto/hmac"
	"crypto/sha256"
	"tyto/core/tyto"
)

// 计算字符串的hmac-sha256值
func HmacSha256Sum(ctx tyto.Context, key string, format string, a ...interface{}) string {
	h := hmac.New(sha256.New, []byte(key))
	return hashSum(ctx, h, format, a...)
}

func HmacSha256ToBase64(ctx tyto.Context, key string, format string, a ...interface{}) string {
	h := hmac.New(sha256.New, []byte(key))
	return hashSumToBase64(ctx, h, format, a...)
}
