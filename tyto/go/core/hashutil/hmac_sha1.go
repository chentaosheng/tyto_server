package hashutil

import (
	"crypto/hmac"
	"crypto/sha1"
	"tyto/core/tyto"
)

// 计算字符串的hmac-sha1值
func HmacSha1Sum(ctx tyto.Context, key string, format string, a ...interface{}) string {
	h := hmac.New(sha1.New, []byte(key))
	return hashSum(ctx, h, format, a...)
}

func HmacSha1ToBase64(ctx tyto.Context, key string, format string, a ...interface{}) string {
	h := hmac.New(sha1.New, []byte(key))
	return hashSumToBase64(ctx, h, format, a...)
}
